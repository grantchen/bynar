package repository

import (
	"database/sql"
	"fmt"
)

type accountManagerRepository struct {
	db *sql.DB
}

type PermissionInfo struct {
	Id             int
	OrganizationId int `field:"organization_id"`
	Archived       int
	Status         int
	Suspended      int
	MStatus        int `field:"mstatus"`
	Enterprise     int
	SecretName     string `field:"secret_name"`
}

// CheckPermission implements AccountManagerRepository
func (a *accountManagerRepository) CheckPermission(accountID int, organizationID int) (*PermissionInfo, bool, error) {
	queryCheck := `SELECT accounts.id, accounts.organization_id, organization.archived, tenants.status,
					tenants_management.status as mstatus,tenants_management.suspended, tenants.enterprise, tenants.secret_name
		FROM accounts
		INNER JOIN organizations ON organizations.id = accounts.id,
		INNER JOIN tenants_management ON tenants_management.organizations_id = organizations.id
		INNER JOIN tenants ON tenants.id = tenants_management.tenants_id
		WHERE accounts.id = ? AND organizations.id = ?
	`
	stmt, err := a.db.Prepare(queryCheck)
	if err != nil {
		return nil, false, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, queryCheck)
	}

	defer stmt.Close()

	rows, err := stmt.Query(accountID, organizationID)
	if err != nil {
		return nil, false, fmt.Errorf("query: [%w], sql string: [%s]", err, queryCheck)
	}

	count := 0
	for rows.Next() {
		count += 1
		permissionInfo := &PermissionInfo{Id: -1, OrganizationId: -1, Archived: -1, Status: -1, Suspended: -1, MStatus: -1,
			Enterprise: -1, SecretName: ""}
		rows.Scan(
			&permissionInfo.Id, &permissionInfo.OrganizationId, &permissionInfo.Archived,
			&permissionInfo.Status, &permissionInfo.MStatus, &permissionInfo.Suspended,
			&permissionInfo.Enterprise, &permissionInfo.SecretName,
		)
		return permissionInfo, true, nil
	}
	defer rows.Close()

	return nil, false, fmt.Errorf("something when wrong when parse result, count rows: %d", count)
}

func NewAccountManagerRepository(db *sql.DB) AccountManagerRepository {
	return &accountManagerRepository{
		db: db,
	}
}
