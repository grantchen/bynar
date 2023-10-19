package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"log"
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
		return nil, fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeBeginTransaction))
	}
	defer tx.Rollback()
	isCommit := true
	fieldsCombinationValidating := []string{"status", "general_product_posting_group_id", "general_business_posting_group_id"}
	for _, field := range fieldsCombinationValidating {
		seenMap := make(map[int]bool)
		for _, gr := range grList {
			if gr[field] != nil {
				status, _ := gr.GetValInt("status")
				value, _ := gr.GetValInt(field)
				// Check if the value is already in the map
				if seenMap[value] && status == 1 {
					// If there is the same value, handle it accordingly.
					isCommit = false
					resp.IO.Result = -1
					resp.IO.Message = fmt.Sprintf("%s: %s: %d", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), value)
					resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
					break
				} else {
					seenMap[value] = true
				}
			}
		}
	}
	if isCommit == true {
		for _, gr := range grList {
			if err := u.handle(tx, gr); err != nil {
				log.Println("Err", err)

				resp.IO.Result = -1
				resp.IO.Message += err.Error() + "\n"
				resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
				break
			}
			resp.Changes = append(resp.Changes, gr)
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
		}
	}
	if isCommit == true {
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, errors.ErrCodeCommitTransaction), err)
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
			return i18n.ErrMsgToI18n(err, u.language)
		}
		err = gr.ValidateOnPositiveNumber(repository.GeneralPostingSetupFieldNames)
		if err != nil {
			return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodePositiveNumber))
		}
		generalPostingSetup, _ := model.ParseGridRow(gr)
		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return err
		}
		status, _ := gr.GetValInt("status")
		if status == 1 {
			for _, field := range fieldsCombinationValidating {
				ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(gr, []string{field})
				if !ok || err != nil {
					return fmt.Errorf("%s: %s: %s", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), gr[field])
				}
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err = gr.ValidateOnRequired(repository.GeneralPostingSetupFieldNames)
		if err != nil {
			return i18n.ErrMsgToI18n(err, u.language)
		}
		err = gr.ValidateOnPositiveNumber(repository.GeneralPostingSetupFieldNames)
		if err != nil {
			return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodePositiveNumber))
		}
		//id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err != nil {
			return err
		}

		if generalPostingSetup.Archived == 1 {
			return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeArchivedUpdate))
		}

		// merge request data and current
		generalPostingSetup, err = model.ParseWithDefaultValue(gr, *generalPostingSetup)
		if err != nil {
			return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeMergeRequest))
		}

		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeInvalidCondition))
		}

		logger.Debug("status: ", generalPostingSetup.Status, "check: ", generalPostingSetup.Status == 1)
		status, _ := gr.GetValInt("status")
		if status == 1 {
			newGr := gr.MergeWithMap(generalPostingSetup.ToMap())
			logger.Debug("newMap", newGr)
			for _, field := range fieldsCombinationValidating {
				ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(newGr, []string{field})
				if !ok || err != nil {
					return fmt.Errorf("%s: %s: %s", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), gr[field])
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
				return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeArchivedDelete))
			}
			err = u.tgGeneralPostingSetupSimpleRepository.Delete(tx, gr)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	default:
		return fmt.Errorf("%s: %s", i18n.Localize(u.language, errors.ErrCodeUndefinedTowType), gr.GetActionType())
	}

	if err != nil {
		return i18n.ErrMsgToI18n(err, u.language)
	}

	return err
}

func (u *uploadService) checkGeneralPostSetupCondition(gps *model.GeneralPostingSetup) error {

	if gps.Archived != 0 && gps.Archived != 1 {
		logger.Debug("gps: ", gps.Archived)
		return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeArchivedNotValid))
	}

	if gps.Status != 0 && gps.Status != 1 {
		return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeStatusNotValid))
	}

	if gps.Status == 1 && gps.Status == gps.Archived {
		return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeSameArchivedStatus))
	}

	return nil
}
