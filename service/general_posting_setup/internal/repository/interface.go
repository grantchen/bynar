package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/model"
)

type GeneralPostingSetupRepository interface {
	IsContainCombination(tx *sql.Tx, status int, generalProductPostingGroupID int, generalBussinessPostingGroupID int) (bool, error)
	GetGeneralPostingSetup(generalPostingSetupID int) (*model.GeneralPostingSetup, error)
	GetGeneralPostingSetupAsMap(generalPostingSetupID int) (map[string]interface{}, error)
}
