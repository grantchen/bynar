package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"github.com/sirupsen/logrus"
)

type treegridService struct {
	db          *sql.DB
	userService service.UserService
}

func newTreeGridService(db *sql.DB, accountID int, organizationUuid string) treegrid.TreeGridService {

	logger.Debug("accountID:", accountID)

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepository(db, "users", repository.UserFieldNames,
		100)

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		logrus.Error(err)
	}
	appConfig := config.NewLocalConfig()
	accountDB, err := sql_db.NewConnection(appConfig.GetAccountManagementConnection())
	if err != nil {
		logrus.Error(err)
	}
	var oid int
	accountDB.QueryRow("SELECT organizations.id FROM organizations WHERE organization_uuid = ?", organizationUuid).Scan(&oid)
	userService := service.NewUserService(db, accountDB, oid, authProvider, simpleOrganizationRepository)

	return &treegridService{
		db:          db,
		userService: *userService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo) treegrid.TreeGridService {
		return newTreeGridService(db, accountID, organizationUuid)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.userService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.userService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.userService.Handle(req)
}
