package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
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
	resp := &treegrid.PostResponse{Changes: []map[string]interface{}{}}
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
	isCommit := true
	for _, gr := range grList {
		if err = u.handle(tx, gr); err != nil {
			log.Println("Err", err)
			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
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

func (u *uploadService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error

	fieldsCombinationValidating := []string{"status", "general_product_posting_group_id", "general_business_posting_group_id"}
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = gr.ValidateOnRequired(repository.GeneralPostingSetupFieldNames)
		if err != nil {
			return i18n.SimpleTranslation(u.language, "RequiredFieldsBlank", nil)
		}
		err = gr.ValidateOnNotNegativeNumber(repository.GeneralPostingSetupFieldNames, u.language)
		if err != nil {
			return err
		}
		generalPostingSetup, _ := model.ParseGridRow(gr)
		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return i18n.SimpleTranslation(u.language, "", err)
		}
		status, _ := gr.GetValInt("status")
		if status == 1 {
			for _, field := range fieldsCombinationValidating {
				ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
				if !ok || err != nil {
					templateData := map[string]string{
						"Field": field,
					}
					return i18n.ParametersTranslation(u.language, "ValueDuplicated", templateData)
				}
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err = gr.ValidateOnRequired(repository.GeneralPostingSetupFieldNames)
		if err != nil {
			return i18n.SimpleTranslation(u.language, "RequiredFieldsBlank", nil)
		}
		err = gr.ValidateOnNotNegativeNumber(repository.GeneralPostingSetupFieldNames, u.language)
		if err != nil {
			return err
		}
		//id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err != nil {
			return i18n.SimpleTranslation(u.language, "", err)
		}

		if generalPostingSetup.Archived == 1 {
			return i18n.SimpleTranslation(u.language, "ArchivedUpdate", nil)
		}

		// merge request data and current
		generalPostingSetup, err = model.ParseWithDefaultValue(gr, *generalPostingSetup)
		if err != nil {
			return i18n.SimpleTranslation(u.language, "MergeFail", nil)
		}

		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return i18n.SimpleTranslation(u.language, "InvalidCondition", nil)
		}

		logger.Debug("status: ", generalPostingSetup.Status, "check: ", generalPostingSetup.Status == 1)
		status, _ := gr.GetValInt("status")
		if status == 1 {
			newGr := gr.MergeWithMap(generalPostingSetup.ToMap())
			logger.Debug("newMap", newGr)
			for _, field := range fieldsCombinationValidating {
				ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(tx, newGr, []string{field})
				if !ok || err != nil {
					templateData := map[string]string{
						"Field": field,
					}
					return i18n.ParametersTranslation(u.language, "ValueDuplicated", templateData)
				}
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		//id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err == nil {
			if generalPostingSetup.Archived == 1 {
				return i18n.SimpleTranslation(u.language, "ArchivedDelete", nil)
			}
			err = u.tgGeneralPostingSetupSimpleRepository.Delete(tx, gr)
			if err != nil {
				return i18n.SimpleTranslation(u.language, "", err)
			}
		} else {
			return nil
		}
	default:
		return err
	}

	if err != nil {
		return i18n.SimpleTranslation(u.language, "", err)
	}

	return err
}

func (u *uploadService) checkGeneralPostSetupCondition(gps *model.GeneralPostingSetup) error {

	if gps.Archived != 0 && gps.Archived != 1 {
		logger.Debug("gps: ", gps.Archived)
		return i18n.SimpleTranslation(u.language, "NotValidArchived", nil)
	}

	if gps.Status != 0 && gps.Status != 1 {
		return i18n.SimpleTranslation(u.language, "NotValidStatus", nil)
	}

	if gps.Status == 1 && gps.Status == gps.Archived {
		return i18n.SimpleTranslation(u.language, "StatusAndArchivedSame", nil)
	}

	return nil
}
