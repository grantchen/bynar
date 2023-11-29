package model

import (
	"encoding/json"
)

type Warehouses struct {
	ID                             int `json:"id,string"`
	Archived                       int `json:"archived,string"`
	Status                         int `json:"status,string"`
	GeneralProductPostingGroupID   int `json:"general_product_posting_group_id,string"`
	GeneralBussinessPostingGroupID int `json:"general_business_posting_group_id,string"`
}

func defaultWarehouses() *Warehouses {
	return &Warehouses{
		Status:                         1,
		Archived:                       0,
		GeneralProductPostingGroupID:   0,
		GeneralBussinessPostingGroupID: 0,
	}
}

func (g Warehouses) ToMap() map[string]interface{} {
	jsonData, _ := json.Marshal(g)

	var m map[string]interface{}
	json.Unmarshal(jsonData, &m)
	return m
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
