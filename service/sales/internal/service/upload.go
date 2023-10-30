package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"log"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
)

type UploadService struct {
	db                              *sql.DB
	updateGRSaleRepository          treegrid.SimpleGridRowRepository
	updateGRSaleRepositoryWithChild treegrid.GridRowRepositoryWithChild
	language                        string
	accountID                       int
	approvalService                 service.ApprovalService
	docSvc                          service.DocumentService
	saleRep                         repository.SaleRepository
	unitRep                         pkg_repository.UnitRepository
	currencyRep                     pkg_repository.CurrencyRepository
	inventoryRep                    pkg_repository.InventoryRepository
	boundFlowRep                    pkg_repository.BoundFlowRepository
}

func NewUploadService(db *sql.DB,
	updateGRSaleRepository treegrid.SimpleGridRowRepository,
	updateGRSaleRepositoryWithChild treegrid.GridRowRepositoryWithChild,
	language string,
	accountID int,
	approvalService service.ApprovalService,
	docSvc service.DocumentService,
	saleRep repository.SaleRepository,
	unitRep pkg_repository.UnitRepository,
	currencyRep pkg_repository.CurrencyRepository,
	inventoryRep pkg_repository.InventoryRepository,
) *UploadService {
	return &UploadService{
		db:                              db,
		updateGRSaleRepository:          updateGRSaleRepository,
		updateGRSaleRepositoryWithChild: updateGRSaleRepositoryWithChild,
		language:                        language,
		accountID:                       accountID,
		approvalService:                 approvalService,
		docSvc:                          docSvc,
		saleRep:                         saleRep,
		unitRep:                         unitRep,
		currencyRep:                     currencyRep,
		inventoryRep:                    inventoryRep,
	}
}

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	b, _ := json.Marshal(req)
	logger.Debug("request: ", string(b))
	trList, err := treegrid.ParseRequestUpload(req, u.updateGRSaleRepositoryWithChild)

	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeBeginTransaction))
	}
	defer tx.Rollback()

	m := make(map[string]interface{})
	for _, tr := range trList.MainRows() {
		if err := u.handle(tx, tr); err != nil {
			log.Println("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += i18n.ErrMsgToI18n(err, u.language).Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(tr.Fields))
			break
		}
		resp.Changes = append(resp.Changes, tr.Fields)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Fields))
		resp.Changes = append(resp.Changes, m)

		for k := range tr.Items {
			resp.Changes = append(resp.Changes, tr.Items[k])
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Items[k]))
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, errors.ErrCodeCommitTransaction), err)
	}

	return resp, nil
}

func (s *UploadService) handle(tx *sql.Tx, tr *treegrid.MainRow) error {
	// Check Approval Order
	ok, err := s.approvalService.Check(tr, s.accountID, s.language)
	if err != nil {
		return fmt.Errorf("check order: [%w], transfer id: %s", err, tr.IDString())
	}

	if !ok {
		return fmt.Errorf("invalid approval order: [%w]: status: %d", errors.ErrForbiddenAction, tr.Status())
	}

	if err := s.save(tx, tr); err != nil {
		return err
	}

	//if err := s.updateGRSaleRepositoryWithChild.Save(tx, tr); err != nil {
	//	return fmt.Errorf("transfer svc save '%s': [%w]", tr.IDString(), err)
	//}

	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := s.GetSaleTx(tx, tr.Fields.GetID())
		if err != nil {
			return fmt.Errorf("get sale service: [%w]", err)
		}

		if err := s.HandleSale(tx, entity); err != nil {
			return fmt.Errorf("handle sale: [%w]", err)
		}

		if err := s.docSvc.Handle(tx, entity.ID, entity.DocumentID, entity.DocumentNo); err != nil {
			return fmt.Errorf("handle document: [%w], modelID: %d, docID: %d, docNo: '%s'", err, entity.ID, entity.DocumentID, entity.DocumentNo)
		}
	}

	return nil
}

func (s *UploadService) save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := s.saveSale(tx, tr); err != nil {
		return fmt.Errorf("%s %s: [%w]",
			i18n.Localize(s.language, errors.ErrCodeSave),
			i18n.Localize(s.language, errors.ErrCodeSale),
			i18n.ErrMsgToI18n(err, s.language))
	}

	if err := s.saveSaleLine(tx, tr, tr.Fields.GetID()); err != nil {
		return fmt.Errorf("%s %s: [%w]",
			i18n.Localize(s.language, errors.ErrCodeSave),
			i18n.Localize(s.language, errors.ErrCodeSaleLine),
			i18n.ErrMsgToI18n(err, s.language))
	}

	return nil
}

func (s *UploadService) saveSale(tx *sql.Tx, tr *treegrid.MainRow) error {
	fieldsValidating := []string{"code"}

	var err error
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = tr.Fields.ValidateOnRequiredAll(repository.SaleFieldNames)
		if err != nil {
			return err
		}

		for _, field := range fieldsValidating {
			ok, err := s.updateGRSaleRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				return fmt.Errorf("duplicate, %s", field)
			}
		}
	case treegrid.GridRowActionChanged:
		err = tr.Fields.ValidateOnRequired(repository.SaleFieldNames)
		if err != nil {
			return err
		}

		for _, field := range fieldsValidating {
			ok, err := s.updateGRSaleRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				return fmt.Errorf("duplicate, %s", field)
			}
		}
	case treegrid.GridRowActionDeleted:
		// ignore id start with CR
		idStr := tr.Fields.GetIDStr()
		if !strings.HasPrefix(idStr, "CR") {
			stmt, err := tx.Prepare("DELETE FROM sale_lines WHERE parent_id = ?")
			if err != nil {
				return err
			}

			defer stmt.Close()

			_, err = stmt.Exec(idStr)
			if err != nil {
				return err
			}
		}
	}

	return s.updateGRSaleRepositoryWithChild.SaveMainRow(tx, tr)
}

func (s *UploadService) saveSaleLine(tx *sql.Tx, tr *treegrid.MainRow, parentID interface{}) error {
	for _, item := range tr.Items {
		logger.Debug("save group line: ", tr, "parentID: ", parentID)

		var err error
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err = item.ValidateOnRequiredAll(map[string][]string{"user_id": repository.SaleLineFieldNames["user_id"]})
			if err != nil {
				return err
			}

			err = s.updateGRSaleRepositoryWithChild.SaveLineAdd(tx, item)
			if err != nil {
				return fmt.Errorf("add child user groups line error: [%w]", err)
			}
		case treegrid.GridRowActionChanged:
			// DO NOTHING WITH ACTION UPDATE, NOT ALLOW UPDATE LINES TABLE
			return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeNoAllowToUpdateChildLine))
		case treegrid.GridRowActionDeleted:
			logger.Debug("delete child")

			// re-assign sale_lines id
			item["id"] = item.GetID()
			err = s.updateGRSaleRepositoryWithChild.SaveLineDelete(tx, item)
			if err != nil {
				return fmt.Errorf("delete child user group line error: [%w]", err)
			}
		default:
			return fmt.Errorf("undefined row type: %s", tr.Fields.GetActionType())

		}
	}
	return nil
}

func (s *UploadService) GetSaleTx(tx *sql.Tx, id interface{}) (*models.Sale, error) {
	return s.saleRep.GetSale(tx, id)
}

func (s *UploadService) HandleSale(tx *sql.Tx, m *models.Sale) error {
	// update quantity
	lines, err := s.saleRep.GetSaleLines(tx, m.ID)
	if err != nil {
		return fmt.Errorf("get lines: [%w]", err)
	}

	disc, err := s.currencyRep.GetDiscount(m.DiscountID)
	if err != nil {
		logger.Debug("get procurement discount")

		// return fmt.Errorf("get discount: [%w], id %d", err, m.InvoiceDiscountID)
	}
	m.Discount = disc

	currency, err := s.currencyRep.GetCurrency(m.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency: [%w]", err)
	}
	logger.Debug("currency.ExchangeRate", currency.ExchangeRate)

	m.CurrencyValue = currency.ExchangeRate

	for _, v := range lines {
		if err := s.HandleLine(tx, m, v); err != nil {
			return fmt.Errorf("handle line: [%w], id: %d", err, v.ID)
		}

		m.TotalDiscount += v.TotalDiscount
		m.TotalDiscountLcy += v.TotalDiscountLcy

		m.TotalVat += v.TotalVat
		m.TotalVatLcy += v.TotalVatLcy

		m.SubtotalExclusiveVat += v.SubtotalExclusiveVat
		m.SubtotalExclusiveVatLcy += v.SubtotalExclusiveVatLcy

		m.TotalExclusiveVat += v.TotalExclusiveVat
		m.TotalExclusiveVatLcy += v.TotalExclusiveVatLcy

		m.TotalInclusiveVat += v.TotalInclusiveVat
		m.TotalInclusiveVatLcy += v.TotalInclusiveVatLcy

		if err := s.saleRep.SaveSaleLine(tx, v); err != nil {
			return fmt.Errorf("save procurement line: [%w]", err)
		}

		currInv, newInv, err := s.handleInventory(tx, v)
		if err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

		if err := s.handleBoundFlows(tx, m, v, 0, currInv, newInv); err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

	}

	return s.saleRep.SaveSale(tx, m)
}

func (s *UploadService) HandleLine(tx *sql.Tx, pr *models.Sale, l *models.SaleLine) (err error) {
	// calc quantity
	unit, err := s.unitRep.Get(l.ItemUnitID)
	if err != nil {
		return fmt.Errorf("get unit: [%w], id %d", err, l.ItemUnitID)
	}
	logger.Debug("Unit id, value", l.ItemUnitID, unit.Value)

	l.ItemUnitValue = unit.Value
	l.Quantity = l.InputQuantity * unit.Value

	// TODO: add parent discount
	// calc discount
	disc, err := s.currencyRep.GetDiscount(l.DiscountID)
	if err != nil {
		logger.Debug("get discount by id: [%w], id %d", err, l.DiscountID)

		// TODO: handle error
		// return fmt.Errorf("get discount by id: [%w], id %d", err, l.DiscountID)
	}

	lineDiscount, err := disc.Calculate(l.SubtotalExclusiveVat)
	if err != nil {
		return fmt.Errorf("calculate line discount: [%w]", err)
	}
	parentDiscount, err := pr.Discount.Calculate(l.SubtotalExclusiveVat)
	if err != nil {
		return fmt.Errorf("calculate parent discount: [%w]", err)
	}
	l.TotalDiscount = lineDiscount + parentDiscount

	l.TotalExclusiveVat = l.SubtotalExclusiveVat - l.TotalDiscount

	// calc vat
	vat, err := s.currencyRep.GetVat(l.VatID)
	if err != nil {
		logger.Debug("get vat: [%w], id %d", err, l.VatID)

		// TODO: handle error
		// return fmt.Errorf("get vat: [%w], id %d", err, l.VatID)
	}

	if l.TotalVat, err = vat.Calculate(l.TotalExclusiveVat); err != nil {
		return fmt.Errorf("calculate vat: [%w]", err)
	}
	l.TotalInclusiveVat = l.TotalExclusiveVat + l.TotalVat

	l.SubtotalExclusiveVatLcy = l.SubtotalExclusiveVat * pr.CurrencyValue
	l.TotalDiscountLcy = l.TotalDiscount * pr.CurrencyValue
	l.TotalExclusiveVatLcy = l.TotalExclusiveVat * pr.CurrencyValue
	l.TotalVatLcy = l.TotalVat * pr.CurrencyValue
	l.TotalInclusiveVatLcy = l.TotalInclusiveVat * pr.CurrencyValue

	return nil
}

func (s *UploadService) handleInventory(tx *sql.Tx, l *models.SaleLine) (currInventory, newInventory models.Inventory, err error) {
	currInventory, err = s.inventoryRep.GetInventory(tx, l.ItemID, l.LocationID)
	if err != nil {
		err = fmt.Errorf("get inventory: [%w], itemID: %d, locationID: %d", err, l.ItemID, l.LocationID)

		return
	}

	newInventory = currInventory

	if currInventory.Quantity == 0 {
		logger.Debug("quantity = 0")

		return
	}

	if currInventory.Quantity < l.Quantity {
		err = fmt.Errorf("outbound quantity is more than available. Outbound: %.2f, available: %.2f", l.Quantity, currInventory.Quantity)

		return
	}

	cost := newInventory.Value / newInventory.Quantity

	// save in inventory only 2 and 3 type
	if l.ItemType != 2 && l.ItemType != 3 {
		return
	}

	newInventory.Value -= l.Quantity * cost
	newInventory.ValueFIFO -= l.Quantity * cost
	newInventory.Quantity -= l.Quantity

	if err = s.inventoryRep.UpdateInventory(tx, newInventory); err != nil {
		err = fmt.Errorf("update inventory: [%w]", err)

		return
	}

	return
}

func (s *UploadService) handleBoundFlows(tx *sql.Tx, pr *models.Sale, l *models.SaleLine, moduleID int, currInv, newInv models.Inventory) (err error) {
	inFlow := models.InboundFlow{
		ModuleID:         moduleID,
		ParentID:         pr.ID,
		PostingDate:      pr.PostingDate,
		ItemID:           currInv.ItemID,
		LocationID:       currInv.LocationID,
		Quantity:         currInv.Quantity,
		Value:            currInv.Value,
		OutboundQuantity: currInv.Quantity - newInv.Quantity,
		OutboundValue:    currInv.Value - newInv.Value,
	}

	if inFlow.Quantity == 0 {
		inFlow.Status = 1
	}

	if err = s.boundFlowRep.SaveInboundFlow(tx, inFlow); err != nil {
		return fmt.Errorf("save inbound flow: [%w]", err)
	}

	outFlow := models.OutboundFlow{
		ModuleID:      moduleID,
		ParentID:      pr.ID,
		TransactionNo: pr.TransactionNo,
		PostingDate:   pr.PostingDate,
		LocationID:    currInv.LocationID,
		ItemID:        currInv.ItemID,
		Quantity:      currInv.Quantity - newInv.Quantity,
		ValueAvco:     currInv.Value - newInv.Value,
		ValueFifo:     currInv.ValueFIFO - newInv.ValueFIFO,
	}

	if err = s.boundFlowRep.SaveOutboundFlow(tx, outFlow); err != nil {
		return fmt.Errorf("save outbound flow: [%w]", err)
	}

	return
}
