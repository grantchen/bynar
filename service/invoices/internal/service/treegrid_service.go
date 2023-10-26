package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"log"
)

type TreeGridService struct {
	db                      *sql.DB
	invoiceSimpleRepository treegrid.SimpleGridRowRepository
	accountID               int
	language                string
}

func NewTreeGridService(db *sql.DB, invoiceSimpleRepository treegrid.SimpleGridRowRepository, accountID int, language string) (*TreeGridService, error) {
	return &TreeGridService{
		db:                      db,
		invoiceSimpleRepository: invoiceSimpleRepository,
		accountID:               accountID,
		language:                language,
	}, nil
}

// GetPageCount implements TreeGridService
func (u *TreeGridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := u.invoiceSimpleRepository.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements TreeGridService
func (u *TreeGridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.invoiceSimpleRepository.GetPageData(tr)
}

// Handle implements TreeGridService
func (u *TreeGridService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{Changes: []map[string]interface{}{}}
	// Create new transaction
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}
	isCommit := true
	// If no errors occurred, commit the transaction
	for _, gr := range grList {
		if err = u.handle(tx, gr); err != nil {
			log.Println("Err", err)
			isCommit = false
			resp.IO.Result = -1
			resp.IO.Message += i18n.ErrMsgToI18n(err, u.language).Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			isCommit = false
			break
		}
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}
	if isCommit == true {
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit transaction: [%w]", err)
		}
	}
	return resp, nil
}

func (s *TreeGridService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error

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
		ok, err1 := s.invoiceSimpleRepository.ValidateOnIntegrity(tx, gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("%s, duplicate", "invoice_no")
		}
		err = s.invoiceSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		// Support operations that are not "update"
		_, ok := gr.GetValInt("id")
		if !ok {
			return nil
		}

		err1 := gr.ValidateOnRequired(repository.InvoiceFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 = s.invoiceSimpleRepository.ValidateOnIntegrity(tx, gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("%s, duplicate", "invoice_no")
		}
		err = s.invoiceSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.invoiceSimpleRepository.Delete(tx, gr)
	default:
		return err
	}

	if err != nil {
		return err
	}

	return err
}
