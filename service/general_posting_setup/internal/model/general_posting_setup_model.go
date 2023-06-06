package model

import (
	"encoding/json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type GeneralPostingSetup struct {
	ID                             int `json:"id,string"`
	Archived                       int `json:"archived,string"`
	Status                         int `json:"status,string"`
	GeneralProductPostingGroupID   int `json:"general_product_posting_group_id,string"`
	GeneralBussinessPostingGroupID int `json:"general_business_posting_group_id,string"`
}

// type generalPostingSetupDTO struct {
// 	ID                             int `json:"id"`
// 	Archived                       int `json:"archived"`
// 	Status                         int `json:"status"`
// 	GeneralProductPostingGroupID   int `json:"general_product_posting_group_id"`
// 	GeneralBussinessPostingGroupID int `json:"general_business_posting_group_id"`
// }

func defaultGeneralPostingSetup() *GeneralPostingSetup {
	return &GeneralPostingSetup{
		Status:                         1,
		Archived:                       0,
		GeneralProductPostingGroupID:   0,
		GeneralBussinessPostingGroupID: 0,
	}
}

func ParseGridRow(gr treegrid.GridRow) (*GeneralPostingSetup, error) {
	return ParseWithDefaultValue(gr, *defaultGeneralPostingSetup())
}

func (g GeneralPostingSetup) ToMap() map[string]interface{} {
	jsonData, _ := json.Marshal(g)

	var m map[string]interface{}
	json.Unmarshal(jsonData, &m)
	return m
}

func ParseWithDefaultValue(gr treegrid.GridRow, defaultValue GeneralPostingSetup) (*GeneralPostingSetup, error) {
	result := &GeneralPostingSetup{}
	*result = defaultValue

	jsonData, err := json.Marshal(gr)
	logger.Debug(result)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(jsonData, &result)
	return result, nil
}

func ParseFromMapStr(input map[string]string) (*GeneralPostingSetup, error) {
	result := defaultGeneralPostingSetup()

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(jsonData, &result)
	return result, nil
}
