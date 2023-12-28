package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type ProcurementSvc struct {
	db                             *sql.DB
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	procRep                        repository.ProcurementRepository
	unitRep                        repository.UnitRepository
	currencyRep                    repository.CurrencyRepository
	inventoryRep                   repository.InventoryRepository
	boundFlowRep                   repository.BoundFlowRepository
	language                       string
}

func NewProcurementSvc(db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild, procRep repository.ProcurementRepository, unitRep repository.UnitRepository,
	currencyRep repository.CurrencyRepository, inventoryRep repository.InventoryRepository, language string) ProcurementsService {
	return &ProcurementSvc{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
		procRep:                        procRep,
		unitRep:                        unitRep,
		currencyRep:                    currencyRep,
		inventoryRep:                   inventoryRep,
		language:                       language,
	}
}

// GetPageCount implements PaymentsService
func (u *ProcurementSvc) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return u.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

// GetPageData implements PaymentsService
func (u *ProcurementSvc) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.gridRowDataRepositoryWithChild.GetPageData(tr)
}

func (u *ProcurementSvc) GetTx(tx *sql.Tx, id interface{}) (*models.Procurement, error) {
	return u.procRep.GetProcurement(tx, id)
}

func (u *ProcurementSvc) Save(tx *sql.Tx, m *models.Procurement) error {
	return u.procRep.SaveProcurement(tx, m)
}

func (u *ProcurementSvc) Handle(tx *sql.Tx, m *models.Procurement) error {
	// update quantity
	lines, err := u.procRep.GetProcurementLines(tx, m.ID)
	if err != nil {
		return fmt.Errorf("get lines: [%w]", err)
	}

	disc, err := u.currencyRep.GetDiscount(m.InvoiceDiscountID)
	if err != nil {
		logger.Debug("get procurement discount")

		// return fmt.Errorf("get discount: [%w], id %d", err, m.InvoiceDiscountID)
	}
	m.Discount = disc
	// get currency
	currency, err := u.currencyRep.GetCurrency(m.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency: [%w]", err)
	}
	logger.Debug("currency.ExchangeRate", currency.ExchangeRate)

	m.CurrencyValue = currency.ExchangeRate

	// handle lines
	for _, v := range lines {
		if err := u.HandleLine(tx, m, v); err != nil {
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

		if err := u.procRep.SaveProcurementLine(tx, v); err != nil {
			return fmt.Errorf("save procurement line: [%w]", err)
		}

		cost, err := u.handleInventory(tx, v)
		if err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

		if err := u.handleInboundFlow(tx, m, v, cost); err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

	}

	return u.procRep.SaveProcurement(tx, m)
}

func (u *ProcurementSvc) HandleLine(_ *sql.Tx, pr *models.Procurement, l *models.ProcurementLine) (err error) {
	// calc quantity
	unit, err := u.unitRep.Get(l.ItemUnitID)
	if err != nil {
		return fmt.Errorf("get unit: [%w], id %d", err, l.ItemUnitID)
	}
	logger.Debug("Unit id, value", l.ItemUnitID, unit.Value)

	l.ItemUnitValue = unit.Value
	l.Quantity = l.InputQuantity * unit.Value
	// calc discount
	disc, err := u.currencyRep.GetDiscount(l.DiscountID)
	if err != nil {
		logger.Debug("get discount by id: [%w], id %d", err, l.DiscountID)

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
	vat, err := u.currencyRep.GetVat(l.VatID)
	if err != nil {
		logger.Debug("get vat: [%w], id %d", err, l.VatID)
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

func (u *ProcurementSvc) handleInventory(tx *sql.Tx, l *models.ProcurementLine) (cost float32, err error) {
	inv, err := u.inventoryRep.GetInventory(tx, l.ItemID, l.LocationID)
	if err != nil {
		return 0, fmt.Errorf("get inventory: [%w], itemID: %d, locationID: %d", err, l.ItemID, l.LocationID)
	}

	if inv.Quantity == 0 {
		logger.Debug("quantity = 0")

		return
	}
	cost = inv.Value / inv.Quantity
	// save in inventory only 2 and 3 type
	if l.ItemType != 2 && l.ItemType != 3 {
		return
	}

	return cost, u.inventoryRep.AddValues(tx, l.ItemID, l.LocationID, l.Quantity, l.TotalInclusiveVatLcy)
}

func (u *ProcurementSvc) handleInboundFlow(tx *sql.Tx, pr *models.Procurement, l *models.ProcurementLine, cost float32) (err error) {
	inFlow := models.InboundFlow{
		PostingDate:      pr.PostingDate,
		ItemID:           l.ItemID,
		ParentID:         l.ParentID,
		LocationID:       l.LocationID,
		Quantity:         l.Quantity,
		Value:            l.TotalInclusiveVat * pr.CurrencyValue,
		OutboundQuantity: l.Quantity,
		OutboundValue:    l.Quantity * cost,
	}

	err = u.boundFlowRep.SaveInboundFlow(tx, inFlow)

	return
}

// validate item_id
func (u *ProcurementSvc) validateItemID(tx *sql.Tx, item treegrid.GridRow) error {
	id, ok := item["item_id"]
	if !ok {
		return nil
	}
	query := `SELECT 1 FROM items WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return i18n.TranslationI18n(u.language, "ItemNotExist", map[string]string{
				"ItemId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(u.language, err)
	}

	return nil
}

func (u *ProcurementSvc) ValidateParams(db *sql.DB, item treegrid.GridRow) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil
	}

	err = u.validateItemID(tx, item)
	if err != nil {
		return err
	}
	unitID, ok := item["item_unit_id"]
	if !ok {
		return nil
	}

	query := `SELECT value FROM units WHERE id = ?`
	var unitValue float64
	if err := tx.QueryRow(query, unitID).Scan(&unitValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("unit not found with item_unit_id: %s", unitID)
		}

		return fmt.Errorf("query row: [%w], query: %s", err, query)
	}
	return nil
}
