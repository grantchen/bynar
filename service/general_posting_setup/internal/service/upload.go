package service

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type uploadService struct {
	db                                    *sql.DB
	tgGeneralPostingSetupSimpleRepository treegrid.SimpleGridRowRepository
	generalPostingSetupRepository         repository.GeneralPostingSetupRepository
	language                              string
}

func NewUploadService(db *sql.DB,
	tgGeneralPostingSetupSimpleRepository treegrid.SimpleGridRowRepository,
	generalPostingSetupRepository repository.GeneralPostingSetupRepository,
	language string,

) UploadService {
	return &uploadService{
		db:                                    db,
		tgGeneralPostingSetupSimpleRepository: tgGeneralPostingSetupSimpleRepository,
		generalPostingSetupRepository:         generalPostingSetupRepository,
		language:                              language,
	}
}

// Handle implements UploadService
func (u *uploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
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

func (u *uploadService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error

	fieldsCombinationValidating := []string{"status", "general_product_posting_group_id", "general_business_posting_group_id"}
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = gr.ValidateOnRequired(repository.GeneralPostingSetupFieldNames, u.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnNotNegativeNumber(repository.GeneralPostingSetupFieldNames, u.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLength(repository.GeneralPostingSetupFieldNames, 100, u.language)
		if err != nil {
			return err
		}
		generalPostingSetup, _ := model.ParseGridRow(gr)
		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return i18n.TranslationErrorToI18n(u.language, err)
		}
		status, _ := gr.GetValInt("status")
		if status == 1 {
			field := "code"
			ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}

			ok, err = u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(tx, gr, fieldsCombinationValidating)
			if !ok || err != nil {
				templateData := map[string]string{
					"Fields": "general_product_posting_group_id, general_business_posting_group_id",
				}
				return i18n.TranslationI18n(u.language, "FieldsCombinationDuplicated", templateData)
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err = gr.ValidateOnRequired(repository.GeneralPostingSetupFieldNames, u.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnNotNegativeNumber(repository.GeneralPostingSetupFieldNames, u.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLength(repository.GeneralPostingSetupFieldNames, 100, u.language)
		if err != nil {
			return err
		}
		//id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err != nil {
			return i18n.TranslationErrorToI18n(u.language, err)
		}

		if generalPostingSetup.Archived == 1 {
			return i18n.TranslationI18n(u.language, "ArchivedUpdate", map[string]string{})
		}

		// merge request data and current
		generalPostingSetup, err = model.ParseWithDefaultValue(gr, *generalPostingSetup)
		if err != nil {
			return i18n.TranslationI18n(u.language, "MergeFail", map[string]string{})
		}

		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return i18n.TranslationI18n(u.language, "InvalidCondition", map[string]string{})
		}

		if generalPostingSetup.Status == 1 {
			newGr := gr.MergeWithMap(generalPostingSetup.ToMap())
			field := "code"
			ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(tx, newGr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}

			ok, err = u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(tx, newGr, fieldsCombinationValidating)
			if !ok || err != nil {
				templateData := map[string]string{
					"Fields": "general_product_posting_group_id, general_business_posting_group_id",
				}
				return i18n.TranslationI18n(u.language, "FieldsCombinationDuplicated", templateData)
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		//id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err == nil {
			if generalPostingSetup.Archived == 1 {
				return i18n.TranslationI18n(u.language, "ArchivedDelete", map[string]string{})
			}
			err = u.tgGeneralPostingSetupSimpleRepository.Delete(tx, gr)
			if err != nil {
				return i18n.TranslationErrorToI18n(u.language, err)
			}
		} else {
			return err
		}
	default:
		return err
	}

	if err != nil {
		return i18n.TranslationErrorToI18n(u.language, err)
	}

	return err
}

func (u *uploadService) checkGeneralPostSetupCondition(gps *model.GeneralPostingSetup) error {

	if gps.Archived != 0 && gps.Archived != 1 {
		logger.Debug("gps: ", gps.Archived)
		return i18n.TranslationI18n(u.language, "NotValidArchived", map[string]string{})
	}

	if gps.Status != 0 && gps.Status != 1 {
		return i18n.TranslationI18n(u.language, "NotValidStatus", map[string]string{})
	}

	if gps.Status == 1 && gps.Status == gps.Archived {
		return i18n.TranslationI18n(u.language, "StatusAndArchivedSame", map[string]string{})
	}

	return nil
}
