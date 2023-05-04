package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type units struct {
	conn *sql.DB
}

func NewUnitRepository(conn *sql.DB) UnitRepository {
	return &units{conn: conn}
}

func (r *units) Get(id interface{}) (m models.Unit, err error) {
	query := `
	SELECT value
	FROM units
	WHERE id = ?
	`

	err = r.conn.QueryRow(query, id).Scan(&m.Value)

	return
}
