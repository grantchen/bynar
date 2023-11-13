package service

import (
	"database/sql"
	"fmt"

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
}

func NewProcurementSvc(db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild, procRep repository.ProcurementRepository, unitRep repository.UnitRepository,
	currencyRep repository.CurrencyRepository, inventoryRep repository.InventoryRepository) ProcurementsService {
	return &ProcurementSvc{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
		procRep:                        procRep,
		unitRep:                        unitRep,
		currencyRep:                    currencyRep,
		inventoryRep:                   inventoryRep,
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

func (s *ProcurementSvc) GetTx(tx *sql.Tx, id interface{}) (*models.Procurement, error) {
	return s.procRep.GetProcurement(tx, id)
}

func (s *ProcurementSvc) Save(tx *sql.Tx, m *models.Procurement) error {
	return s.procRep.SaveProcurement(tx, m)
}

func (s *ProcurementSvc) Handle(tx *sql.Tx, m *models.Procurement) error {
	// update quantity
	lines, err := s.procRep.GetProcurementLines(tx, m.ID)
	if err != nil {
		return fmt.Errorf("get lines: [%w]", err)
	}

	disc, err := s.currencyRep.GetDiscount(m.InvoiceDiscountID)
	if err != nil {
		logger.Debug("get procurement discount")

		// return fmt.Errorf("get discount: [%w], id %d", err, m.InvoiceDiscountID)
	}
	m.Discount = disc
	// get currency
	currency, err := s.currencyRep.GetCurrency(m.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency: [%w]", err)
	}
	logger.Debug("currency.ExchangeRate", currency.ExchangeRate)

	m.CurrencyValue = currency.ExchangeRate

	// handle lines
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

		if err := s.procRep.SaveProcurementLine(tx, v); err != nil {
			return fmt.Errorf("save procurement line: [%w]", err)
		}

		cost, err := s.handleInventory(tx, v)
		if err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

		if err := s.handleInboundFlow(tx, m, v, cost); err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

	}

	return s.procRep.SaveProcurement(tx, m)
}

func (s *ProcurementSvc) HandleLine(tx *sql.Tx, pr *models.Procurement, l *models.ProcurementLine) (err error) {
	// calc quantity
	unit, err := s.unitRep.Get(l.ItemUnitID)
	if err != nil {
		return fmt.Errorf("get unit: [%w], id %d", err, l.ItemUnitID)
	}
	logger.Debug("Unit id, value", l.ItemUnitID, unit.Value)

	l.ItemUnitValue = unit.Value
	l.Quantity = l.InputQuantity * unit.Value
	// calc discount
	disc, err := s.currencyRep.GetDiscount(l.DiscountID)
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
	vat, err := s.currencyRep.GetVat(l.VatID)
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

func (s *ProcurementSvc) handleInventory(tx *sql.Tx, l *models.ProcurementLine) (cost float32, err error) {
	inv, err := s.inventoryRep.GetInventory(tx, l.ItemID, l.LocationID)
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

	return cost, s.inventoryRep.AddValues(tx, l.ItemID, l.LocationID, l.Quantity, l.TotalInclusiveVatLcy)
}

func (s *ProcurementSvc) handleInboundFlow(tx *sql.Tx, pr *models.Procurement, l *models.ProcurementLine, cost float32) (err error) {
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

	err = s.boundFlowRep.SaveInboundFlow(tx, inFlow)

	return
}
