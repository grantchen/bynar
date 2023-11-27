package model

type WarehouseInfo struct {
	Id                 string `json:"id,omitempty"`
	WName              string `json:"w_name,omitempty"`
	CostingMethods     string `json:"costing_methods,omitempty"`
	Uuid               string `json:"uuid,omitempty"`
	ProductUuid        string `json:"product_uuid,omitempty"`
	ProductDescription string `json:"product_description,omitempty"`
	ProductBarcode     string `json:"product_barcode,omitempty"`
	Quantity           string `json:"quantity,omitempty"`
	AvgCost            string `json:"avg_cost,omitempty"`
	Name               string `json:"GridName,omitempty"`
}

type Grid struct {
	Changes []*Change `json:"Changes,omitempty"`
}

type Change struct {
	ID   string    `json:"id,omitempty"`
	Data *DataItem `json:"warehouse_descriptionSuggest,omitempty"`
}

type DataItem struct {
	Items []*WarehouseInfo `json:"Items,omitempty"`
}
