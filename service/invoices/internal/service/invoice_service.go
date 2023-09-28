package service

import (
	"database/sql"
)

type invoiceService struct {
	db *sql.DB
}

func NewInvoiceService(db *sql.DB) InvoiceService {
	return &invoiceService{db: db}
}
