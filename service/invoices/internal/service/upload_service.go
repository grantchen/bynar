package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService struct {
	db                      *sql.DB
	invoiceSimpleRepository treegrid.SimpleGridRowRepository
	accountID               int
}

func NewUploadService(db *sql.DB,
	invoiceSimpleRepository treegrid.SimpleGridRowRepository,
	accountID int,
) (*UploadService, error) {
	return &UploadService{
		db:                      db,
		invoiceSimpleRepository: invoiceSimpleRepository,
		accountID:               accountID,
	}, nil
}

// GetPageCount implements treegrid.UploadService
func (u *UploadService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := u.invoiceSimpleRepository.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.UploadService
func (u *UploadService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.invoiceSimpleRepository.GetPageData(tr)
}

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{Changes: []map[string]interface{}{}}
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

	fieldsValidating := []string{"invoice_no"}

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		// Assigning values to other fields
		gr["account_id"] = s.accountID
		err1 := gr.ValidateOnRequiredAll(repository.InvoiceFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := s.invoiceSimpleRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%v], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = s.invoiceSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.InvoiceFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := s.invoiceSimpleRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%w], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = s.invoiceSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.invoiceSimpleRepository.Delete(tx, gr)

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
