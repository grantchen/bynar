package repository

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type generalSettingUploadRepository struct {
}

// GetGeneralPostingSetup implements UploadRepository
func (g *generalSettingUploadRepository) GetGeneralPostingSetup(tx *sql.Tx, generalPostingSetupID int) (*model.GeneralPostingSetup, error) {
	query := QuerySelectWithoutJoin + " WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, query)
	}
	defer stmt.Close()

	rows, err := stmt.Query(generalPostingSetupID)
	if err != nil {
		return nil, fmt.Errorf("query error: [%w], sql string: [%s]", err, query)
	}
	rowVals, err := utils.NewRowVals(rows)

	if err != nil {
		return nil, fmt.Errorf("parse new row error: [%w], sql string: [%s]", err, query)
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

// IsContainCombination implements UploadRepository
func (*generalSettingUploadRepository) IsContainCombination(tx *sql.Tx, status int, generalProductPostingGroupID int, generalBussinessPostingGroupID int) (bool, error) {
	panic("unimplemented")
}

func NewPostingSetupRepository() GeneralPostingSetupRepository {
	return &generalSettingUploadRepository{}
}
