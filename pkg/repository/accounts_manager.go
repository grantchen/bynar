package repository

import (
	"database/sql"
	"fmt"
)

type accountManagerRepository struct {
	db *sql.DB
}

type PermissionInfo struct {
}

// CheckPermission implements AccountManagerRepository
func (a *accountManagerRepository) CheckPermission(accountID int, organizationID int) (bool, error) {
	queryCheck := `SELECT accounts.id, accounts.organizations_id, organization.archived, tenants.status,
					tenants_management.status,tenants_management.suspended
		FROM accounts
		INNER JOIN organizations ON organizations.id = accounts.id,
		INNER JOIN tenants_management ON tenants_management.organizations_id = organizations.id
		INNER JOIN tenants ON tenants.id = tenants_management.tenants_id
		WHERE accounts.id = ? AND organizations.id = ?
	`
	stmt, err := a.db.Prepare(queryCheck)
	if err != nil {
		return false, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, queryCheck)
	}

	defer stmt.Close()

	rows, err := stmt.Query(accountID, organizationID)
	if err != nil {
		return false, fmt.Errorf("query: [%w], sql string: [%s]", err, queryCheck)
	}

	count := 0
	for rows.Next() {
		count += 1
		rows.Scan()
	}
	defer rows.Close()

}

func NewAccountManagerRepository(db *sql.DB) AccountManagerRepository {
	return &accountManagerRepository{
		db: db,
	}
}
