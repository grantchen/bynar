package service

import (
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/repository"
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

	fieldsCombinationValidating := []string{"code"}
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.WarehousesFieldNames, u.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnNotNegativeNumber(repository.WarehousesFieldNames, u.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLength(repository.WarehousesFieldNames, 100, u.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsCombinationValidating {
			ok, err := u.tgWarehousesSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}
		}
		err = u.tgWarehousesSimpleRepository.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.WarehousesFieldNames, u.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnNotNegativeNumber(repository.WarehousesFieldNames, u.language)
		if err != nil {
			return err
		}
		err = gr.ValidateOnLimitLength(repository.WarehousesFieldNames, 100, u.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsCombinationValidating {
			ok, err := u.tgWarehousesSimpleRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}
		}
		_, ok := gr.GetValInt("id")
		if ok {
			err = u.tgWarehousesSimpleRepository.Update(tx, gr)
		}
	case treegrid.GridRowActionDeleted:
		err = u.tgWarehousesSimpleRepository.Delete(tx, gr)
		//id := gr.GetIDInt()
		var warehouses *model.Warehouses
		warehouses, err = u.warehousesRepository.GetWarehouses(gr.GetIDInt())
		if err == nil {
			// check archived status
			if warehouses.Archived == 1 {
				return i18n.TranslationI18n(u.language, "ArchivedDelete", map[string]string{})
			}
			err = u.tgWarehousesSimpleRepository.Delete(tx, gr)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	default:
		return i18n.TranslationErrorToI18n(u.language, err)
	}

	if err != nil {
		return i18n.TranslationErrorToI18n(u.language, err)
	}

	return err
}

// checkGeneralPostSetupCondition
func (u *uploadService) checkGeneralPostSetupCondition(gps *model.Warehouses) error {

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
