package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type uploadService struct {
	db                                    *sql.DB
	tgGeneralPostingSetupSimpleRepository treegrid.SimpleGridRowRepository
	generalPostingSetupRepository         repository.GeneralPostingSetupRepository
}

func NewUploadService(db *sql.DB,
	tgGeneralPostingSetupSimpleRepository treegrid.SimpleGridRowRepository,
	generalPostingSetupRepository repository.GeneralPostingSetupRepository,
) UploadService {
	return &uploadService{
		db:                                    db,
		tgGeneralPostingSetupSimpleRepository: tgGeneralPostingSetupSimpleRepository,
		generalPostingSetupRepository:         generalPostingSetupRepository,
	}
}

// Handle implements UploadService
func (u *uploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
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
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			break
		}
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}

	return resp, nil
}

func (u *uploadService) handle(gr treegrid.GridRow) error {
	var err error
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	fieldsCombinationValidating := []string{"status", "general_product_posting_group_id", "general_business_posting_group_id"}
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:

		err = gr.ValidateOnRequired(repository.GeneralPostingSetupFieldNames)
		if err != nil {
			return err
		}
		generalPostingSetup, _ := model.ParseGridRow(gr)
		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return err
		}

		if generalPostingSetup.Status == 1 {

			ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(gr, fieldsCombinationValidating)
			if !ok || err != nil {
				return fmt.Errorf("validate duplicate when add: [%v], field: %s, status = 1", err, strings.Join(fieldsCombinationValidating, ", "))
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:

		id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err != nil {
			return err
		}

		if generalPostingSetup.Archived == 1 {
			return fmt.Errorf("archived = 1, can not update on id: [%d]", id)
		}

		// merge request data and current
		generalPostingSetup, err = model.ParseWithDefaultValue(gr, *generalPostingSetup)
		if err != nil {
			return fmt.Errorf("merge with current data fail: [%w], id = [%d]", err, id)
		}

		err = u.checkGeneralPostSetupCondition(generalPostingSetup)
		if err != nil {
			return fmt.Errorf("invalid condition when update data: [%w]", err)
		}

		logger.Debug("status: ", generalPostingSetup.Status, "check: ", generalPostingSetup.Status == 1)
		if generalPostingSetup.Status == 1 {
			newGr := gr.MergeWithMap(generalPostingSetup.ToMap())
			logger.Debug("newMap", newGr)
			ok, err := u.tgGeneralPostingSetupSimpleRepository.ValidateOnIntegrity(newGr, fieldsCombinationValidating)
			if !ok || err != nil {
				return fmt.Errorf("validate duplicate when update: [%v], field: %s, status = 1", err, strings.Join(fieldsCombinationValidating, ", "))
			}
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		id := gr.GetIDInt()
		var generalPostingSetup *model.GeneralPostingSetup
		generalPostingSetup, err = u.generalPostingSetupRepository.GetGeneralPostingSetup(gr.GetIDInt())
		if err != nil {
			return err
		}

		if generalPostingSetup.Archived == 1 {
			return fmt.Errorf("archived = 1, can not update on id: [%d]", id)
		}
		err = u.tgGeneralPostingSetupSimpleRepository.Delete(tx, gr)

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

func (u *uploadService) checkGeneralPostSetupCondition(gps *model.GeneralPostingSetup) error {

	if gps.Archived != 0 && gps.Archived != 1 {
		logger.Debug("gps: ", gps.Archived)
		return fmt.Errorf("not valid archived value")
	}

	if gps.Status != 0 && gps.Status != 1 {
		return fmt.Errorf("not valid status value")
	}

	if gps.Status == 1 && gps.Status == gps.Archived {
		return fmt.Errorf("status and archived is same value 1")
	}

	return nil
}
