package service

import (
	"database/sql"
	"fmt"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	sql_connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/scope"
)

type accountManagerRepository struct {
	db                       *sql.DB
	accountManagerRepository repository.AccountManagerRepository
	secretmanager            secretsmanager.SecretsManager
}

func NewAccountManagerService(
	db *sql.DB,
	accaccountManagerRepository repository.AccountManagerRepository,
	secretmanager secretsmanager.SecretsManager,
) AccountManagerService {
	return &accountManagerRepository{
		db:                       db,
		accountManagerRepository: accaccountManagerRepository,
		secretmanager:            secretmanager,
	}
}

// CheckPermission implements AccountManagerService
func (a *accountManagerRepository) CheckPermission(token string) (*repository.PermissionInfo, bool, error) {
	requestScope, err := scope.ResolveFromToken(token)
	if err != nil {
		return nil, false, fmt.Errorf("check permission error: [%w]", err)
	}

	permission, ok, err := a.accountManagerRepository.CheckPermission(requestScope.AccountID, requestScope.OrganizationID)

	if !ok || err != nil {
		return permission, false, err
	}

	// check permission
	ok = permission.Archived == 0 && permission.Status == 1 && permission.MStatus == 1 && permission.Suspended == 0
	return permission, ok, err
}

// GetNewStringConnection implements AccountManagerService
func (a *accountManagerRepository) GetNewStringConnection(token string, permission *repository.PermissionInfo) (string, error) {
	value, err := a.secretmanager.GetString(permission.SecretName)

	if err != nil {
		logger.Debug(err)
		return "", err

	}
	connectionString := sql_connection.JSON2DatabaseConnection(value)

	if permission.Enterprise == 1 {
		return connectionString, nil
	}

	connectionString = sql_connection.ChangeDatabaseConnectionSchema(connectionString, strconv.Itoa(permission.OrganizationId))
	return connectionString, nil
}
