package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/repository"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService struct {
	db                   *sql.DB
	siteService          SiteService
	siteSimpleRepository treegrid.SimpleGridRowRepository
	language             string
}

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

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	for _, gr := range grList {
		if err = u.handle(tx, gr); err != nil {
			log.Println("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			break
			//rollback
			//return resp, err
		}
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: [%w]", err)
	}

	return resp, nil
}

func (s *UploadService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	fieldsValidating := []string{"code"}
	positiveFieldsValidating := []string{"subsidiaries_uuid", "address_uuid", "contact_uuid", "responsibility_center_uuid"}

	var err error
	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = gr.ValidateOnRequiredAll(repository.SiteFieldNames)
		if err != nil {
			return i18n.SimpleTranslation(s.language, "RequiredFieldsBlank", nil)
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
				return i18n.ParametersTranslation(s.language, "ValueDuplicated", templateData)
			}
		}
		err = s.siteSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err = gr.ValidateOnRequired(repository.SiteFieldNames)
		if err != nil {
			return i18n.SimpleTranslation(s.language, "RequiredFieldsBlank", nil)
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
				return i18n.ParametersTranslation(s.language, "ValueDuplicated", templateData)
			}
		}
		err = s.siteSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.siteSimpleRepository.Delete(tx, gr)

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	return i18n.SimpleTranslation(s.language, "", err)
}
