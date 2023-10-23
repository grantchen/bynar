package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/repository"
	"log"
)

type uploadService struct {
	db                           *sql.DB
	tgWarehousesSimpleRepository treegrid.SimpleGridRowRepository
	warehousesRepository         repository.WarehousesRepository
	language                     string
}

func NewUploadService(db *sql.DB,
	tgWarehousesSimpleRepository treegrid.SimpleGridRowRepository,
	warehousesRepository repository.WarehousesRepository,
	language string,

) UploadService {
	return &uploadService{
		db:                           db,
		tgWarehousesSimpleRepository: tgWarehousesSimpleRepository,
		warehousesRepository:         warehousesRepository,
		language:                     language,
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

	fieldsCombinationValidating := []string{"code"}
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.WarehousesFieldNames)
		if err1 != nil {
			return i18n.ErrMsgToI18n(err1, u.language)
		}
		for _, field := range fieldsCombinationValidating {
			ok, err := u.tgWarehousesSimpleRepository.ValidateOnIntegrity(gr, []string{field})
			if !ok || err != nil {
				return fmt.Errorf("%s: %s: %s", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), gr[field])
			}
		}
		err = u.tgWarehousesSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.WarehousesFieldNames)
		if err1 != nil {
			return err1
		}
		for _, field := range fieldsCombinationValidating {
			ok, err := u.tgWarehousesSimpleRepository.ValidateOnIntegrity(gr, []string{field})
			if !ok || err != nil {
				return fmt.Errorf("%s: %s: %s", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), gr[field])
			}
		}
		err = u.tgWarehousesSimpleRepository.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = u.tgWarehousesSimpleRepository.Delete(tx, gr)
		//id := gr.GetIDInt()
		var warehouses *model.Warehouses
		warehouses, err = u.warehousesRepository.GetWarehouses(gr.GetIDInt())
		if err == nil {
			if warehouses.Archived == 1 {
				return fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeArchivedDelete))
			}
			err = u.tgWarehousesSimpleRepository.Delete(tx, gr)
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

func (u *uploadService) checkGeneralPostSetupCondition(gps *model.Warehouses) error {

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
