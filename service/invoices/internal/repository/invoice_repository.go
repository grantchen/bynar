package repository

import "database/sql"

type invoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}
