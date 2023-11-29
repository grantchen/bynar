package repository

// WarehousesFieldNames is a map of field names to sql column names
var WarehousesFieldNames = map[string][]string{
	"id":                         {"id"},
	"code":                       {"code"},
	"description":                {"description"},
	"transaction_code":           {"transaction_code"},
	"site_uuid":                  {"site_uuid"},
	"address_uuid":               {"address_uuid"},
	"contact_uuid":               {"contact_uuid"},
	"responsibility_center_uuid": {"responsibility_center_uuid"},
}
