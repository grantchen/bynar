package service

import (
	"context"
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db             *sql.DB
	paymentService service.PaymentService
	uploadService  *service.UploadService
}

func newTreeGridService(db *sql.DB, language string) treegrid.TreeGridService {
	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"payments",
		"payment_lines",
		repository.PaymentFieldNames,
		repository.PaymentLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "document_id",
			QueryParent:              repository.QueryParent,
			QueryParentCount:         repository.QueryParentCount,
			QueryParentJoins:         repository.QueryParentJoins,
			QueryChild:               repository.QueryChild,
			QueryChildCount:          repository.QueryChildCount,
			QueryChildJoins:          repository.QueryChildJoins,
			QueryChildSuggestion:     repository.QueryChildSuggestion,
			ChildJoinFieldWithParent: "parent_id",
			ParentIdField:            "id",
		},
	)
	userGroupService := service.NewUserGroupService(db, gridRowDataRepositoryWithChild)

	grUserGroupDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "user_groups",
		"user_group_lines",
		repository.UserGroupFieldNames,
		repository.UserGroupLineFieldUploadNames)

	grUserRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"users",
		repository.UserUploadNames,
		1, // arbitrary
	)

	grUserGroupRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"user_groups",
		repository.UserGroupFieldNames,
		1, // arbitrary
	)

	uploadService := service.NewUploadService(db, grUserGroupRepository, grUserGroupDataUploadRepositoryWithChild, grUserRepository, language)
	return &treegridService{
		db:               db,
		userGroupService: userGroupService,
		uploadService:    uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, AccountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (s *treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	return s.userGroupService.GetCellSuggestion(req)
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.userGroupService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.userGroupService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
