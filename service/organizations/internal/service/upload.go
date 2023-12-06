package service

import (
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// UploadService is the service for upload
type UploadService struct {
	db                           *sql.DB
	organizationService          OrganizationService
	organizationSimpleRepository treegrid.SimpleGridRowRepository
	userID                       int
	language                     string
}

// NewUploadService create new instance of UploadService
func NewUploadService(db *sql.DB,
	organizationService OrganizationService,
	organizationSimpleRepository treegrid.SimpleGridRowRepository,
	userID int,
	language string,
) (*UploadService, error) {
	return &UploadService{
		db:                           db,
		organizationService:          organizationService,
		organizationSimpleRepository: organizationSimpleRepository,
		userID:                       userID,
		language:                     language,
	}, nil
}

// Handle hanldes the upload request
func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	resp := treegrid.HandleSingleRows(grList, func(gr treegrid.GridRow) error {
		err = utils.WithTransaction(u.db, func(tx *sql.Tx) error {
			return u.handle(tx, gr)
		})
		return i18n.TranslationErrorToI18n(u.language, err)
	})

	return resp, nil
}

// handle handles upload request of single row
func (s *UploadService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error
	fieldsValidating := []string{"code"}

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		var parentId int
		stmt, err1 := s.db.Prepare(`SELECT parent_id FROM user_group_lines WHERE user_id = ?`)
		if err1 != nil {
			return errors.NewUnknownError("prepare sql error", errors.ErrCode).WithInternalCause(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(s.userID).Scan(&parentId)
		if err != nil {
			return i18n.TranslationI18n(s.language, "NoUserGroupLineFound", map[string]string{})
		}

		gr["user_group_int"] = parentId

		err1 = gr.ValidateOnRequiredAll(repository.OrganizationFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnLimitLength(repository.OrganizationFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := s.organizationSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
			}
		}
		err = s.organizationSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.OrganizationFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnLimitLength(repository.OrganizationFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := s.organizationSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
			}
		}
		err = s.organizationSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.organizationSimpleRepository.Delete(tx, gr)

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	if err != nil {
		return i18n.TranslationErrorToI18n(s.language, err)
	}

	return nil
}
