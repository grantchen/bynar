package model

import (
	"encoding/json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type Warehouses struct {
	ID                             int `json:"id,string"`
	Archived                       int `json:"archived,string"`
	Status                         int `json:"status,string"`
	GeneralProductPostingGroupID   int `json:"general_product_posting_group_id,string"`
	GeneralBussinessPostingGroupID int `json:"general_business_posting_group_id,string"`
}

// type warehousesDTO struct {
// 	ID                             int `json:"id"`
// 	Archived                       int `json:"archived"`
// 	Status                         int `json:"status"`
// 	GeneralProductPostingGroupID   int `json:"general_product_posting_group_id"`
// 	GeneralBussinessPostingGroupID int `json:"general_business_posting_group_id"`
// }

func defaultWarehouses() *Warehouses {
	return &Warehouses{
		Status:                         1,
		Archived:                       0,
		GeneralProductPostingGroupID:   0,
		GeneralBussinessPostingGroupID: 0,
	}
}

func ParseGridRow(gr treegrid.GridRow) (*Warehouses, error) {
	return ParseWithDefaultValue(gr, *defaultWarehouses())
}

func (g Warehouses) ToMap() map[string]interface{} {
	jsonData, _ := json.Marshal(g)

	var m map[string]interface{}
	json.Unmarshal(jsonData, &m)
	return m
}

func ParseWithDefaultValue(gr treegrid.GridRow, defaultValue Warehouses) (*Warehouses, error) {
	result := &Warehouses{}
	*result = defaultValue

	jsonData, err := json.Marshal(gr)
	logger.Debug(result)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(jsonData, &result)
	return result, nil
}

func ParseFromMapStr(input map[string]string) (*Warehouses, error) {
	result := defaultWarehouses()

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(jsonData, &result)
	return result, nil
}
