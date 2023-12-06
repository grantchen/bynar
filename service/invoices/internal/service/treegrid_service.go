package service

import (
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
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

func (s *TreeGridService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error

	fieldsValidating := []string{"invoice_no"}
	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		// Assigning values to other fields
		gr["account_id"] = s.accountID
		err1 := gr.ValidateOnRequiredAll(repository.InvoiceFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnLimitLength(repository.InvoiceFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLengthToFloat(repository.InvoiceFieldNamesFloat, s.language)
		if err != nil {
			return err
		}
		ok, err1 := s.invoiceSimpleRepository.ValidateOnIntegrity(tx, gr, fieldsValidating)
		if !ok || err1 != nil {
			templateData := map[string]string{
				"Field": "invoice_no",
			}
			return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
		}
		err = s.invoiceSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		// Support operations that are not "update"
		_, ok := gr.GetValInt("id")
		if !ok {
			return nil
		}

		err1 := gr.ValidateOnRequired(repository.InvoiceFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnLimitLength(repository.InvoiceFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLengthToFloat(repository.InvoiceFieldNamesFloat, s.language)
		if err != nil {
			return err
		}
		ok, err1 = s.invoiceSimpleRepository.ValidateOnIntegrity(tx, gr, fieldsValidating)
		if !ok || err1 != nil {
			templateData := map[string]string{
				"Field": "invoice_no",
			}
			return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
		}
		err = s.invoiceSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.invoiceSimpleRepository.Delete(tx, gr)
	default:
		return err
	}

	if err != nil {
		return i18n.TranslationErrorToI18n(s.language, err)
	}

	return err
}
