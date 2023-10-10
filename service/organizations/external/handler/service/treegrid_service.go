package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db                  *sql.DB
	organizationService service.OrganizationService
	uploadService       service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.TreeGridService {

	var filterPermissionCondition string
	var userID int

	logger.Debug("accountID:", accountID)
	if accountID != 0 {
		appConfig := config.NewLocalConfig()
		accountDB, err := sql_db.NewConnection(appConfig.GetAccountManagementConnection())
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
		defer stmt.Close()
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

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository, userID)
	return &treegridService{
		db:                  db,
		organizationService: organizationService,
		uploadService:       *uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.organizationService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.organizationService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
