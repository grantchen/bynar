package service

import (
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService struct {
	db                       *sql.DB
	languageService          LanguageService
	languageSimpleRepository treegrid.SimpleGridRowRepository
	language                 string
}

func NewUploadService(db *sql.DB,
	languageService LanguageService,
	languageSimpleRepository treegrid.SimpleGridRowRepository,
	language string,
) (*UploadService, error) {
	return &UploadService{
		db:                       db,
		languageService:          languageService,
		languageSimpleRepository: languageSimpleRepository,
		language:                 language,
	}, nil
}

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

func (s *UploadService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error
	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.LanguageFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = validateGridRow(gr, s.language)
		if err != nil {
			return err
		}
		err = s.languageSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.LanguageFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = validateGridRow(gr, s.language)
		if err != nil {
			return err
		}
		err = s.languageSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = s.languageSimpleRepository.Delete(tx, gr)

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	return i18n.TranslationErrorToI18n(s.language, err)
}

// Common verification logic
func validateGridRow(gr treegrid.GridRow, language string) error {
	positiveFieldsValidating := []string{"number"}
	err := gr.ValidateOnLimitLength(repository.LanguageFieldCountry, 30, language)
	if err != nil {
		return err
	}
	err = gr.ValidateOnLimitLength(repository.LanguageFieldLanguage, 40, language)
	if err != nil {
		return err
	}
	err = gr.ValidateOnLimitLength(repository.LanguageFieldLetters, 10, language)
	if err != nil {
		return err
	}
	err = gr.ValidateOnLimitLengthToFloat(repository.LanguageFieldNamesFloat, language)
	if err != nil {
		return err
	}
	for _, field := range positiveFieldsValidating {
		err = gr.ValidateOnNotNegativeNumber(map[string][]string{field: repository.LanguageFieldNames[field]}, language)
		if err != nil {
			return err
		}
	}

	return nil
}
