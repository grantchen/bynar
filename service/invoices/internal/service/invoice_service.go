package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type invoiceService struct {
	db                      *sql.DB
	simpleInvoiceRepository treegrid.SimpleGridRowRepository
}

// GetPageCount implements InvoiceService
func (o *invoiceService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return o.simpleInvoiceRepository.GetPageCount(tr)
}

// GetPageData implements InvoiceService
func (o *invoiceService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return o.simpleInvoiceRepository.GetPageData(tr)
}

func NewInvoiceService(db *sql.DB, simpleInvoiceService treegrid.SimpleGridRowRepository) InvoiceService {
	return &invoiceService{db: db, simpleInvoiceRepository: simpleInvoiceService}
}
