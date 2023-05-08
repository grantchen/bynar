package repository

import (
	"database/sql"
)

type AccountingPeriods struct {
	conn *sql.DB
}

func NewAccountingPeriods(conn *sql.DB) *AccountingPeriods {
	return &AccountingPeriods{conn: conn}
}
