package service

import (
	"context"
	"database/sql"
	"encoding/json"
	stderr "errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"log"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
)

// UploadService is the service for upload
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

// NewUploadService returns new instance of UploadService
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
	boundFlowRep pkg_repository.BoundFlowRepository,
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
		boundFlowRep:                    boundFlowRep,
	}
}

// Handle handles upload request
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
		return nil, fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	m := make(map[string]interface{})
	var handleErr error
	for _, tr := range trList.MainRows() {
		if handleErr = u.handle(tx, tr); handleErr != nil {
			log.Println("Err", handleErr)

			resp.IO.Result = -1
			resp.IO.Message += handleErr.Error() + "\n"
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

	if handleErr == nil {
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit transaction: [%w]", err)
		}
	}

	return resp, nil
}

// handle handles upload request of single row
func (s *UploadService) handle(tx *sql.Tx, tr *treegrid.MainRow) error {
	// Check Approval Order
	ok, err := s.approvalService.Check(tr, s.accountID, s.language)
	if err != nil {
		return err
	}

	if !ok {
		return i18n.TranslationI18n(s.language, "ForbiddenAction", map[string]string{
			"Message": err.Error(),
		})
	}

	if err := s.save(tx, tr); err != nil {
		return err
	}

	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := s.GetSaleTx(tx, tr.Fields.GetID())
		if err != nil {
			return fmt.Errorf("get sale service: [%w]", err)
		}

		if err := s.HandleSale(tx, entity); err != nil {
			return err
		}

		if err := s.docSvc.Handle(tx, entity.ID, entity.DocumentID, entity.DocumentNo); err != nil {
			return fmt.Errorf("handle document: [%w], modelID: %d, docID: %d, docNo: '%s'", err, entity.ID, entity.DocumentID, entity.DocumentNo)
		}
	}

	return nil
}

// save saves sale and sale lines
func (s *UploadService) save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := s.saveSale(tx, tr); err != nil {
		return i18n.TranslationI18n(s.language, "SaveSale", map[string]string{
			"Message": err.Error(),
		})
	}

	if err := s.saveSaleLine(tx, tr, tr.Fields.GetID()); err != nil {
		return i18n.TranslationI18n(s.language, "SaveSaleLine", map[string]string{
			"Message": err.Error(),
		})
	}

	return nil
}

// saveSale saves sale
func (s *UploadService) saveSale(tx *sql.Tx, tr *treegrid.MainRow) error {
	requiredFieldsMapping := tr.Fields.FilterFieldsMapping(
		repository.SaleFieldNames,
		[]string{
			"document_id",
			"document_no",
			"transaction_no",
			"document_date",
			"posting_date",
			"entry_date",
			"shipment_date",
			"store_id",
			"currency_id",
		})
	positiveFieldsMapping := tr.Fields.FilterFieldsMapping(
		repository.SaleFieldNames,
		[]string{
			"document_id",
			"store_id",
			"currency_id",
		})

	var err error
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = tr.Fields.ValidateOnRequiredAll(requiredFieldsMapping, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnLimitLength(requiredFieldsMapping, 100, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnPositiveNumber(positiveFieldsMapping, s.language)
		if err != nil {
			return err
		}

		err = s.validateStoreID(tx, tr)
		if err != nil {
			return err
		}

		tr.Fields["currency_value"] = 0
		tr.Fields["direct_debit_mandate_id"] = 0
	case treegrid.GridRowActionChanged:
		err = tr.Fields.ValidateOnRequired(requiredFieldsMapping, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnLimitLength(requiredFieldsMapping, 100, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnPositiveNumber(positiveFieldsMapping, s.language)
		if err != nil {
			return err
		}

		err = s.validateStoreID(tx, tr)
		if err != nil {
			return err
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

// saveSaleLine saves sale lines
func (s *UploadService) saveSaleLine(tx *sql.Tx, tr *treegrid.MainRow, parentID interface{}) error {
	requiredFieldsMapping := tr.Fields.FilterFieldsMapping(
		repository.SaleLineFieldNames,
		[]string{
			"item_id",
			"item_unit_id",
		})
	positiveFieldsMapping := tr.Fields.FilterFieldsMapping(
		repository.SaleLineFieldNames,
		[]string{
			"item_id",
			"input_quantity",
			"item_unit_id",
		})

	for _, item := range tr.Items {
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err := item.ValidateOnRequiredAll(requiredFieldsMapping, s.language)
			if err != nil {
				return err
			}

			err = tr.Fields.ValidateOnLimitLength(positiveFieldsMapping, 100, s.language)
			if err != nil {
				return err
			}

			err = item.ValidateOnPositiveNumber(positiveFieldsMapping, s.language)
			if err != nil {
				return err
			}

			// check item_id
			if err = s.validateItemID(tx, item); err != nil {
				return err
			}

			// check item_unit_id
			if err = s.validateItemUnitID(tx, item); err != nil {
				return err
			}

			item["item_unit_value"] = 0
			item["quantity"] = 0
			item["quantity_assign"] = 0
			item["quantity_assigned"] = 0
			item["total_exclusive_vat"] = 0
			item["total_vat"] = 0
			item["subtotal_exclusive_vat_lcy"] = 0
			item["total_discount_lcy"] = 0
			item["total_exclusive_vat_lcy"] = 0
			item["total_vat_lcy"] = 0
			item["total_inclusive_vat_lcy"] = 0

			err = s.updateGRSaleRepositoryWithChild.SaveLineAdd(tx, item)
			if err != nil {
				return fmt.Errorf("add child user groups line error: [%w]", err)
			}
		case treegrid.GridRowActionChanged:
			err := item.ValidateOnRequired(requiredFieldsMapping, s.language)
			if err != nil {
				return err
			}

			err = tr.Fields.ValidateOnLimitLength(positiveFieldsMapping, 100, s.language)
			if err != nil {
				return err
			}

			err = item.ValidateOnPositiveNumber(positiveFieldsMapping, s.language)
			if err != nil {
				return err
			}

			// check item_id
			if err = s.validateItemID(tx, item); err != nil {
				return err
			}

			// check item_unit_id
			if err = s.validateItemUnitID(tx, item); err != nil {
				return err
			}

			return s.updateGRSaleRepositoryWithChild.SaveLineUpdate(tx, item)
		case treegrid.GridRowActionDeleted:
			err := s.updateGRSaleRepositoryWithChild.SaveLineDelete(tx, item)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("undefined row type: %s", tr.Fields.GetActionType())

		}
	}
	return nil
}

// GetSaleTx returns sale by id
func (s *UploadService) GetSaleTx(tx *sql.Tx, id interface{}) (*models.Sale, error) {
	return s.saleRep.GetSale(tx, id)
}

// HandleSale handles sale calculation and saving to db
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
		if stderr.Is(err, sql.ErrNoRows) {
			return i18n.TranslationI18n(s.language, "CurrencyNotExist", map[string]string{
				"CurrencyId": strconv.Itoa(m.CurrencyID),
			})
		}
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

// HandleLine handles sale line calculation
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

// handleInventory handles inventory calculation and saving to db
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

// handleBoundFlows handles bound flows calculation and saving to db
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

// validate store_id
func (s *UploadService) validateStoreID(tx *sql.Tx, tr *treegrid.MainRow) error {
	id, ok := tr.Fields["store_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM stores WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return i18n.TranslationI18n(s.language, "StoreNotExist", map[string]string{
				"StoreId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(s.language, err)
	}

	return nil
}

// validate item_id
func (s *UploadService) validateItemID(tx *sql.Tx, item treegrid.GridRow) error {
	id, ok := item["item_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM items WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return i18n.TranslationI18n(s.language, "ItemNotExist", map[string]string{
				"ItemId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(s.language, err)
	}

	return nil
}

// validate item_unit_id
func (s *UploadService) validateItemUnitID(tx *sql.Tx, item treegrid.GridRow) error {
	id, ok := item["item_unit_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM units WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return i18n.TranslationI18n(s.language, "UnitNotExist", map[string]string{
				"ItemUnitId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(s.language, err)
	}

	return nil
}
