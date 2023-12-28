package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// treegridService implements treegrid.Service
type treegridService struct {
	db                  *sql.DB
	organizationService service.OrganizationService
	uploadService       service.UploadService
}

// newTreeGridService create new treegrid service
func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.Service {
	var filterPermissionCondition string
	var userID int

	if accountID != 0 {
		appConfig := config.NewLocalConfig()
		accountDB, err := sqldb.NewConnection(appConfig.GetAccountManagementConnection())
		if err != nil {
			logrus.Error(err)
			return nil
		}

		querySql := fmt.Sprintf(`SELECT organization_user_id FROM organization_accounts WHERE organization_user_uid = (SELECT uid from accounts WHERE id = ?)`)
		stmt, err := accountDB.Prepare(querySql)
		if err != nil {
			logrus.Errorf("db prepare: [%v], sql string: [%s]", err, querySql)
			return nil
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		err = stmt.QueryRow(accountID).Scan(&userID)
		if err != nil {
			logrus.Errorf("query user id: [%v], sql string: [%s]", err, querySql)
			return nil
		}

		filterPermissionCondition = fmt.Sprintf(repository.QueryPermissionFormat, userID)
	}

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "organizations", repository.OrganizationFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "code",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: filterPermissionCondition,
		})
	organizationService := service.NewOrganizationService(db, simpleOrganizationRepository)

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository, userID, language)
	return &treegridService{
		db:                  db,
		organizationService: organizationService,
		uploadService:       *uploadService,
	}
}

// NewTreeGridServiceFactory create new treegrid service factory
func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (*treegridService) GetCellData(_ context.Context, _ *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.organizationService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.organizationService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
