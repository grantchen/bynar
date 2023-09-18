package service

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
)

type accountManagerRepository struct {
	db                       *sql.DB
	accountManagerRepository repository.AccountManagerRepository
	authProvider             gip.AuthProvider
}

func NewAccountManagerService(
	db *sql.DB,
	accaccountManagerRepository repository.AccountManagerRepository,
	authProvider gip.AuthProvider,
) AccountManagerService {
	return &accountManagerRepository{
		db:                       db,
		accountManagerRepository: accaccountManagerRepository,
		authProvider:             authProvider,
	}
}

// CheckPermission implements AccountManagerService
func (a *accountManagerRepository) CheckPermission(claims *middleware.IdTokenClaims) (*repository.PermissionInfo, bool, error) {

	// TODO:
	permission, ok, err := a.accountManagerRepository.CheckPermission(0, 0)

	if !ok || err != nil {
		return permission, false, err
	}

	// check permission
	ok = permission.Archived == 0 && permission.Status == 1 && permission.MStatus == 1 && permission.Suspended == 0
	return permission, ok, err
}

// GetNewStringConnection implements AccountManagerService
func (a *accountManagerRepository) GetNewStringConnection(tenantUuid, organizationUuid string, permission *repository.PermissionInfo) (string, error) {
	if len(os.Getenv(tenantUuid)) == 0 {
		return "", errors.New("no mysql conn environment of " + tenantUuid)
	}
	envs := strings.Split(os.Getenv(tenantUuid), "/")
	connStr := envs[0] + "/" + organizationUuid
	if len(envs) > 1 {
		connStr += envs[1]
	}
	return connStr, nil
}

// GetRole implements AccountManagerService
func (a *accountManagerRepository) GetRole(accountID int) (map[string]int, error) {
	return a.accountManagerRepository.CheckRole(accountID)
}
