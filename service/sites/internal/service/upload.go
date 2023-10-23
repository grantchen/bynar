package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/repository"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
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

	for _, gr := range grList {
		if err := u.handle(gr); err != nil {
			log.Println("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += i18n.ErrMsgToI18n(err, u.language).Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			break
		}
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}

	return resp, nil
}

func (s *UploadService) handle(gr treegrid.GridRow) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	fieldsValidating := []string{"code"}
	positiveFieldsValidating := []string{"subsidiaries_uuid", "address_uuid", "contact_uuid", "responsibility_center_uuid"}

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = gr.ValidateOnRequiredAll(repository.SiteFieldNames)
		if err != nil {
			return err
		}

		for _, field := range positiveFieldsValidating {
			err = gr.ValidateOnPositiveNumber(map[string][]string{field: repository.SiteFieldNames[field]})
			if err != nil {
				return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodePositiveNumber))
			}
		}

		for _, field := range fieldsValidating {
			ok, err := s.siteSimpleRepository.ValidateOnIntegrity(gr, []string{field})
			if !ok || err != nil {
				return fmt.Errorf("%s: %s: %s", field, i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr[field])
			}
		}
		err = s.siteSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err = gr.ValidateOnRequired(repository.SiteFieldNames)
		if err != nil {
			return err
		}

		for _, field := range positiveFieldsValidating {
			err = gr.ValidateOnPositiveNumber(map[string][]string{field: repository.SiteFieldNames[field]})
			if err != nil {
				return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodePositiveNumber))
			}
		}

		for _, field := range fieldsValidating {
			ok, err := s.siteSimpleRepository.ValidateOnIntegrity(gr, []string{field})
			if !ok || err != nil {
				return fmt.Errorf("%s: %s: %s", field, i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr[field])
			}
		}
		err = s.siteSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.siteSimpleRepository.Delete(tx, gr)

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: [%w]", err)
	}

	return err
}
