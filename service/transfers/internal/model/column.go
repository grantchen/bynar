package model

import "strings"

type Column struct {
	DBName      string
	DBNameShort string
	GridName    string
	IsItem      bool
	IsDate      bool
}

func NewColumn(gridName string) (c Column) {
	c.GridName = gridName
	if val, ok := ItemsFields[gridName]; ok {
		c.IsItem = true
		c.DBName = val
	}

	if val, ok := FieldAliases[gridName]; ok {
		c.DBName = val
	}

	if c.DBName == "" {
		c.DBName = gridName
	}

	nameParts := strings.Split(c.DBName, ".")
	if len(nameParts) > 0 {
		c.DBNameShort = nameParts[len(nameParts)-1]
	}

	return
}

var (
	ItemsFields = map[string]string{
		"item_type":          "item_types.code",
		"item_no":            "items.no",
		"item_name":          "items.description",
		"item_unit":          "units.code",
		"input_quantity":     "transfers_items.input_quantity",
		"item_quantity_unit": "transfers_items.item_quantity_unit",
		"item_quantity":      "transfers_items.item_quantity",
		"item_tempory":       "transfers_items.item_tempory",
		"item_uuid":          "transfers_items.item_uuid",
		"item_unit_uuid":     "transfers_items.item_unit_uuid",
		"item_code":          "items.code",
		"item_barcode":       "items.code_2",
	}

	FieldAliases = map[string]string{
		"document_abbrevation": "documents.document_abbrevation",
		"document_type":        "documents.document_type",
		"document_no":          "transfers.document_no",
		// "document_abbrevation": "transfers.document_abbrevation",
		// "document_type":              "transfers.document_type",
		"store_origin_code":          "stores.code",
		"warehouse_origin_code":      "wh_origin.code",
		"warehouse_destination_code": "wh_destination.code",
		"store_destination_code":     "stores.code",
		"responsibility_center":      `responsibility_center.code`,
		"document_date":              `STR_TO_DATE(document_date,'%m/%d/%Y')`,
		"posting_date":               `STR_TO_DATE(posting_date,'%m/%d/%Y')`,
		"delivery_date":              `STR_TO_DATE(delivery_date,'%m/%d/%Y')`,
		"entry_date":                 `STR_TO_DATE(entry_date,'%m/%d/%Y')`,
	}

	// TransferFields = map[strin]

	FieldAliasesDate = map[string]string{
		"1": " = ",
		"2": " != ",
		"3": " < ",
		"4": " <= ",
		"5": " > ",
		"6": " >= ",
	}

	TransferItemsFields = map[string]bool{
		"item_uuid":    true,
		"item_name":    true,
		"item_code":    true,
		"item_type":    true,
		"item_barcode": true,
		// "item_brand":       true,
		// "item_category":    true,
		// "item_subcategory": true,
		"item_unit":     true,
		"item_quantity": true,
	}

	TransferFields = map[string]bool{
		"transfers.id":                               true,
		"transfers.document_no":                      true,
		"transfers.document_date":                    true,
		"transfers.posting_date":                     true,
		"transfers.entry_date":                       true,
		"transfers.delivery_date":                    true,
		"documents.document_type":                    true,
		"documents.document_abbrevation":             true,
		"stores.code":                                true,
		"wh_origin.code":                             true,
		"wh_destination.code":                        true,
		"ss.code":                                    true,
		"responsibility_center.code":                 true,
		"transfers.document_type_uuid":               true,
		"transfers.store_origin_uuid":                true,
		"transfers.warehouse_origin_uuid":            true,
		"transfers.warehouse_destination_uuid":       true,
		"transfers.responsibility_center_uuid":       true,
		"transfers.warehouseman_destination_approve": true,
		"transfers.has_child":                        true,
	}
)

func IsChildItem(val string) bool {
	return TransferItemsFields[val]
}
