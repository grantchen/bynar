package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/service"
)

type treegridService struct {
	db               *sql.DB
	userGroupService service.UserGroupService
	uploadService    *service.UploadService
}

func newTreeGridService(db *sql.DB) treegrid.TreeGridService {
	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"user_groups",
		"user_group_lines",
		repository.UserGroupFieldNames,
		repository.UserGroupLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "code",
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

	uploadService := service.NewUploadService(db, grUserGroupDataUploadRepositoryWithChild, grUserRepository)
	return &treegridService{
		db:               db,
		userGroupService: userGroupService,
		uploadService:    uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, AccountID int, permissionInfo *treegrid.PermissionInfo) treegrid.TreeGridService {
		return newTreeGridService(db)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(ctx context.Context, tr *treegrid.Treegrid) (float64, error) {
	count, err := s.userGroupService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.userGroupService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(ctx context.Context, req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
