package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type accountManagerRepository struct {
	db *sql.DB
}

type PermissionInfo struct {
	Id               int
	OrganizationId   int `field:"organization_id"`
	TMOrganizationId int `field:"tm_organization_id"`
	Archived         int
	Status           int
	Suspended        int
	MStatus          int `field:"mstatus"`
	Enterprise       int
	SecretName       string `field:"secret_name"`
}

// CheckPermission implements AccountManagerRepository
func (a *accountManagerRepository) CheckPermission(accountID int, organizationID int) (*PermissionInfo, bool, error) {
	queryCheck := `SELECT organizations.id as organization_id, organizations.archived, tenants.status,
						tenants_management.status as mstatus,tenants_management.suspended, tenants.enterprise, tenants.secret_name,
						tenants_management.organization_id as tm_organization_id
						FROM organizations
						INNER JOIN tenants ON tenants.organizations = organizations.id
						INNER JOIN tenants_management ON tenants.id = tenants_management.tenant_id
						WHERE organizations.id = ?
	`
	stmt, err := a.db.Prepare(queryCheck)
	if err != nil {
		return nil, false, fmt.Errorf("db prepare: [%w], sql string: [%s], organization: [%d]", err, queryCheck, organizationID)
	}

	defer stmt.Close()

	rows, err := stmt.Query(organizationID)
	if err != nil {
		return nil, false, fmt.Errorf("query: [%w], sql string: [%s]", err, queryCheck)
	}

	count := 0
	for rows.Next() {
		count += 1
		permissionInfo := &PermissionInfo{Id: -1, OrganizationId: -1, Archived: -1, Status: -1, Suspended: -1, MStatus: -1,
			Enterprise: -1, SecretName: ""}
		rows.Scan(
			&permissionInfo.OrganizationId, &permissionInfo.Archived,
			&permissionInfo.Status, &permissionInfo.MStatus, &permissionInfo.Suspended,
			&permissionInfo.Enterprise, &permissionInfo.SecretName,
			&permissionInfo.TMOrganizationId,
		)
		return permissionInfo, true, nil
	}
	defer rows.Close()

	return nil, false, nil
}

// CheckRole implements AccountManagerRepository
func (a *accountManagerRepository) CheckRole(accountID int) (map[string]int, error) {
	query :=
		`
	SELECT policies.* FROM policies
	INNER JOIN accounts ON accounts.policy_id = policies.id
	WHERE accounts.id = ?
	`
	logger.Debug("query: ", query, accountID)
	stmt, err := a.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s], account: [%d]", err, query, accountID)
	}

	defer stmt.Close()

	rows, err := stmt.Query(accountID)
	if err != nil {
		return nil, fmt.Errorf("query: [%w], sql string: [%s]", err, query)
	}

	rowValues, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row values: [%w], values: [%v]", err, rowValues)
	}

	resultString := make(map[string]string, 0)
	for rows.Next() {
		if err := rowValues.Parse(rows); err != nil {
			return nil, fmt.Errorf("parse rows: [%w]", err)
		}
		resultString = rowValues.StringValues()
	}

	if len(resultString) > 0 {
		result := make(map[string]int, 0)

		for k, v := range resultString {

			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			result[k] = i

		}

		return result, nil

	}

	defer rows.Close()

	return nil, fmt.Errorf("not policies found with accountID: %d", accountID)
}

func NewAccountManagerRepository(db *sql.DB) AccountManagerRepository {
	return &accountManagerRepository{
		db: db,
	}
}
