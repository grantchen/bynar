package service

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
)

type saleService struct {
	saleRep      repository.SaleRepository
	unitRep      pkg_repository.UnitRepository
	currencyRep  pkg_repository.CurrencyRepository
	inventoryRep pkg_repository.InventoryRepository
	boundFlowRep pkg_repository.BoundFlowRepository
}

func NewSaleService(saleRep repository.SaleRepository, unitRep pkg_repository.UnitRepository,
	currencyRep pkg_repository.CurrencyRepository, inventoryRep pkg_repository.InventoryRepository) SaleService {
	return &saleService{
		saleRep:      saleRep,
		unitRep:      unitRep,
		currencyRep:  currencyRep,
		inventoryRep: inventoryRep,
	}
}

func (s *saleService) GetSaleTx(tx *sql.Tx, id interface{}) (*models.Sale, error) {
	return s.saleRep.GetSale(tx, id)
}

func (s *saleService) Handle(tx *sql.Tx, m *models.Sale, moduleID int) error {
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

		if err := s.handleBoundFlows(tx, m, v, moduleID, currInv, newInv); err != nil {
			return fmt.Errorf("handle inventory: [%w], id: %d", err, v.ID)
		}

	}

	return s.saleRep.SaveSale(tx, m)
}

func (s *saleService) HandleLine(tx *sql.Tx, pr *models.Sale, l *models.SaleLine) (err error) {
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

func (s *saleService) handleInventory(tx *sql.Tx, l *models.SaleLine) (currInventory, newInventory models.Inventory, err error) {
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

func (s *saleService) handleBoundFlows(tx *sql.Tx, pr *models.Sale, l *models.SaleLine, moduleID int, currInv, newInv models.Inventory) (err error) {
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
