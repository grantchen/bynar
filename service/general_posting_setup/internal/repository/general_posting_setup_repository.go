package repository

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type generalSettingUploadRepository struct {
	db *sql.DB
}

// GetGeneralPostingSetupAsMap implements GeneralPostingSetupRepository
func (g *generalSettingUploadRepository) GetGeneralPostingSetupAsMap(generalPostingSetupID int) (map[string]interface{}, error) {
	rows, err := g.getGeneralPostingSetup(generalPostingSetupID)
	if err != nil {
		return nil, err
	}
	rowVals, err := utils.NewRowVals(rows)

	if err != nil {
		return nil, fmt.Errorf("parse new row error: [%w]", err)
	}

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return nil, fmt.Errorf("parse rows: [%w]", err)
		}

		return rowVals.Values(), nil
	}

	return nil, fmt.Errorf("not found general posting setup with id: [%d]", generalPostingSetupID)
}

// GetGeneralPostingSetup implements UploadRepository
func (g *generalSettingUploadRepository) GetGeneralPostingSetup(generalPostingSetupID int) (*model.GeneralPostingSetup, error) {
	rows, err := g.getGeneralPostingSetup(generalPostingSetupID)
	if err != nil {
		return nil, err
	}
	rowVals, err := utils.NewRowVals(rows)

	if err != nil {
		return nil, fmt.Errorf("parse new row error: [%w]", err)
	}

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return nil, fmt.Errorf("parse rows: [%w]", err)
		}

		data := rowVals.StringValues()
		generalPostingSetup, err := model.ParseFromMapStr(data)
		if err != nil {
			return nil, fmt.Errorf("parse json: [%w]", err)
		}
		return generalPostingSetup, nil
	}

	return nil, fmt.Errorf("not found general posting setup with id: [%d]", generalPostingSetupID)
}

func (g *generalSettingUploadRepository) getGeneralPostingSetup(generalPostingSetupID int) (*sql.Rows, error) {
	query := QuerySelectWithoutJoin + " WHERE id = ?"
	stmt, err := g.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, query)
	}
	defer stmt.Close()

	rows, err := stmt.Query(generalPostingSetupID)
	if err != nil {
		return nil, fmt.Errorf("query error: [%w], sql string: [%s]", err, query)
	}
	return rows, nil
}

// IsContainCombination implements UploadRepository
func (*generalSettingUploadRepository) IsContainCombination(tx *sql.Tx, status int, generalProductPostingGroupID int, generalBussinessPostingGroupID int) (bool, error) {
	panic("unimplemented")
}

func NewPostingSetupRepository(db *sql.DB) GeneralPostingSetupRepository {
	return &generalSettingUploadRepository{db: db}
}
