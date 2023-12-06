package model

import (
	"encoding/json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type GeneralPostingSetup struct {
	Code                           string `json:"code"`
	Archived                       int    `json:"archived,string"`
	Status                         int    `json:"status,string"`
	GeneralProductPostingGroupID   int    `json:"general_product_posting_group_id,string"`
	GeneralBussinessPostingGroupID int    `json:"general_business_posting_group_id,string"`
}

func defaultGeneralPostingSetup() *GeneralPostingSetup {
	return &GeneralPostingSetup{
		Code:                           "",
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

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ParseFromMapStr(input map[string]string) (*GeneralPostingSetup, error) {
	result := defaultGeneralPostingSetup()

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
