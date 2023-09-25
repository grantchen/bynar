package model

import (
	"database/sql"
)

type Account struct {
	ID         int            `db:"id"`
	Email      sql.NullString `db:"email"`
	FullName   sql.NullString `db:"full_name"`
	Address    sql.NullString `db:"address"`
	Address2   sql.NullString `db:"address_2"`
	Phone      sql.NullString `db:"phone"`
	City       sql.NullString `db:"city"`
	PostalCode sql.NullString `db:"postal_code"`
	Country    sql.NullString `db:"country"`
	State      sql.NullString `db:"state"`
	Status     sql.NullInt64  `db:"status"`
	UID        sql.NullString `db:"uid"`
	OrgID      sql.NullString `db:"org_id"`
	Verified   sql.NullInt64  `db:"verified"`
	Policies   sql.NullString `db:"policies"`
}
