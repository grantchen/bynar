package service

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db          *sql.DB
	userService service.UserService
}

func newTreeGridService(db *sql.DB, accountID int, organizationUuid, language string) treegrid.Service {

	logger.Debug("accountID:", accountID)

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "users", repository.UserFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{MainCol: "email"})

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		logrus.Error(err)
	}

	paymentProvider, err := checkout.NewPaymentClient()
	if err != nil {
		logrus.Error(err)
	}
	appConfig := config.NewLocalConfig()
	accountDB, err := sqldb.NewConnection(appConfig.GetAccountManagementConnection())
	if err != nil {
		logrus.Error(err)
	}
	var oid int
	_ = accountDB.QueryRow("SELECT organizations.id FROM organizations WHERE organization_uuid = ?", organizationUuid).Scan(&oid)
	var customerID string
	_ = accountDB.QueryRow(`SELECT user_payment_gateway_id FROM accounts_cards ac JOIN accounts a on ac.user_id = a.id WHERE ac.user_id = ? AND is_default = ?`, accountID, true).Scan(&customerID)
	userService := service.NewUserService(db, accountDB, oid, customerID, authProvider, paymentProvider, simpleOrganizationRepository, language)

	return &treegridService{
		db:          db,
		userService: *userService,
	}
}

func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, organizationUuid, language)
	}
}

// GetCellData implements treegrid.Service
func (*treegridService) GetCellData(_ context.Context, _ *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.userService.GetPageCount(tr)
	return count, err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.userService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.userService.Handle(req)
}
