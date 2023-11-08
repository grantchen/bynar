package repository

var InvoiceFieldNames = map[string][]string{
	"id":                  {"id"},
	"account_id":          {"account_id"},
	"invoice_date":        {"invoice_date"},
	"invoice_no":          {"invoice_no"},
	"currency":            {"currency"},
	"total":               {"total"},
	"billing_period_date": {"billing_period_date"},
	"provider_id":         {"provider_id"},
	"paid":                {"paid"},
}

var InvoiceFieldNamesFloat = map[string][]string{
	"total": {"total"},
}
