package repository

// OrganizationFieldNames is a map of field names for the organizations table
var OrganizationFieldNames = map[string][]string{
	"id":             {"id"},
	"name":           {"name"},
	"vat_no":         {"vat_number"},
	"state":          {"state"},
	"code":           {"code"},
	"user_group_int": {"user_group_int"},
}
