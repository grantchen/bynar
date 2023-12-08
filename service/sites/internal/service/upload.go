package service

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/repository"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// UploadService is the service for upload
type UploadService struct {
	db                   *sql.DB
	siteService          SiteService
	siteSimpleRepository treegrid.SimpleGridRowRepository
	language             string
}

// NewUploadService create new instance of UploadService
func NewUploadService(db *sql.DB,
	siteService SiteService,
	siteSimpleRepository treegrid.SimpleGridRowRepository,
	language string,
) (*UploadService, error) {
	return &UploadService{
		db:                   db,
		siteService:          siteService,
		siteSimpleRepository: siteSimpleRepository,
		language:             language,
	}, nil
}

// Handle handles upload request
func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	resp := treegrid.HandleSingleTreegridRows(grList, func(gr treegrid.GridRow) error {
		err = utils.WithTransaction(u.db, func(tx *sql.Tx) error {
			return u.handle(tx, gr)
		})
		return i18n.TranslationErrorToI18n(u.language, err)
	})

	return resp, nil
}

// handle handles upload request of single row
func (s *UploadService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	fieldsValidating := []string{"code"}
	positiveFieldsValidating := []string{"subsidiaries_uuid", "address_uuid", "contact_uuid", "responsibility_center_uuid"}

	var err error
	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = gr.ValidateOnRequiredAll(repository.SiteFieldNames, s.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLength(repository.SiteFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range positiveFieldsValidating {
			err = gr.ValidateOnNotNegativeNumber(map[string][]string{field: repository.SiteFieldNames[field]}, s.language)
			if err != nil {
				return err
			}
		}

		for _, field := range fieldsValidating {
			ok, err := s.siteSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
			}
		}
		err = s.siteSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err = gr.ValidateOnRequired(repository.SiteFieldNames, s.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLength(repository.SiteFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range positiveFieldsValidating {
			err = gr.ValidateOnNotNegativeNumber(map[string][]string{field: repository.SiteFieldNames[field]}, s.language)
			if err != nil {
				return err
			}
		}

		for _, field := range fieldsValidating {
			ok, err := s.siteSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
			}
		}
		err = s.siteSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.siteSimpleRepository.Delete(tx, gr)

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	return i18n.TranslationErrorToI18n(s.language, err)
}
