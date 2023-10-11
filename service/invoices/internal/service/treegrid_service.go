package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"log"
	"strings"
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
		return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, errors.ErrCodeBeginTransaction), err)
	}
	defer tx.Rollback()

	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}
	isCommit := true
	seenInvoiceNos := make(map[string]bool)

	for _, gr := range grList {
		if gr["invoice_no"] != nil {
			invoiceNo := gr["invoice_no"].(string)
			//Check if the value is already in the map
			if seenInvoiceNos[invoiceNo] {
				//If there is the same invoice_no, handle it accordingly.
				isCommit = false
				resp.IO.Result = -1
				resp.IO.Message = fmt.Sprintf("invoice_no: %s%s:%s", i18n.Localize(u.language, errors.ErrCodeValueDuplicated), ":", gr["invoice_no"].(string))
				resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			} else {
				seenInvoiceNos[invoiceNo] = true
			}
		}
	}
	// If no errors occurred, commit the transaction
	if isCommit == true {
		for _, gr := range grList {
			if err = u.handle(tx, gr); err != nil {
				log.Println("Err", err)
				isCommit = false
				resp.IO.Result = -1
				resp.IO.Message += err.Error() + "\n"
				resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
				break
			}
			resp.Changes = append(resp.Changes, gr)
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
		}
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, errors.ErrCodeCommitTransaction), err)
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
			return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeRequiredFieldsBlank))
		}
		ok, err1 := s.invoiceSimpleRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("%s: %s: %s", strings.Join(fieldsValidating, ", "), i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["invoice_no"])
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
			return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeRequiredFieldsBlank))
		}
		ok, err1 = s.invoiceSimpleRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("%s %s: %s", strings.Join(fieldsValidating, ", "), i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["invoice_no"])
		}
		err = s.invoiceSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.invoiceSimpleRepository.Delete(tx, gr)
	default:
		return fmt.Errorf("%s: %s", i18n.Localize(s.language, errors.ErrCodeUndefinedTowType), gr.GetActionType())
	}

	if err != nil {
		//Formatted messy string
		contains := strings.Contains(err.Error(), "of range")
		if contains {
			return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeOutRange))
		} else {
			return fmt.Errorf(err.Error())
		}
	}

	return err
}
