package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db              *sql.DB
	languageservice service.languageservice
	uploadService   *service.UploadService
}

func newTreeGridService(db *sql.DB, language string) treegrid.TreeGridService {
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
	languageservice := service.Newlanguageservice(db, gridRowDataRepositoryWithChild)

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
		db:              db,
		languageservice: languageservice,
		uploadService:   uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, AccountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (s *treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	return s.languageservice.GetCellSuggestion(req)
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.languageservice.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.languageservice.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
