package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService struct {
	db                           *sql.DB
	organizationService          OrganizationService
	organizationSimpleRepository treegrid.SimpleGridRowRepository
}

func NewUploadService(db *sql.DB,
	organizationService OrganizationService,
	organizationSimpleRepository treegrid.SimpleGridRowRepository,
) (*UploadService, error) {
	return &UploadService{
		db:                           db,
		organizationService:          organizationService,
		organizationSimpleRepository: organizationSimpleRepository,
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
			resp.IO.Message += err.Error() + "\n"
			break
		}
		resp.Changes = append(resp.Changes, gr)
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

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.OrganizationFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := s.organizationSimpleRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%v], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = s.organizationSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.OrganizationFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := s.organizationSimpleRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%w], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = s.organizationSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.organizationSimpleRepository.Delete(tx, gr)

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
