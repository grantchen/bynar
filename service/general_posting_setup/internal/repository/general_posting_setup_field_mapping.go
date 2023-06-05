package repository

var GeneralPostingSetupFieldNames = map[string][]string{
	"id":                                       {"id"},
	"sales_account":                            {"sales_account"},
	"sales_credit_memo_account":                {"sales_credit_memo_account"},
	"sales_line_discount_account":              {"sales_line_discount_account"},
	"sales_invoice_discount_account":           {"sales_invoice_discount_account"},
	"sales_payment_discount_debit_account":     {"sales_payment_discount_debit_account"},
	"sales_payment_discount_credit_account":    {"sales_payment_discount_credit_account"},
	"code":                                     {"code"},
	"description":                              {"description"},
	"purchase_account":                         {"purchase_account"},
	"sales_prepayments_account":                {"sales_prepayments_account"},
	"purchase_credit_memo_account":             {"purchase_credit_memo_account"},
	"purchase_line_discount_account":           {"purchase_line_discount_account"},
	"purchase_inventory_discount_account":      {"purchase_inventory_discount_account"},
	"purchase_payment_discount_debit_account":  {"purchase_payment_discount_debit_account"},
	"purchase_payment_discount_credit_account": {"purchase_payment_discount_credit_account"},
	"purchase_prepayments_account":             {"purchase_prepayments_account"},
	"cogs_account":                             {"cogs_account"},
	"inventory_adjustment_account":             {"inventory_adjustment_account"},
	"overhead_applied_account":                 {"overhead_applied_account"},
	"purchase_variance_account":                {"purchase_variance_account"},
	"used_ledger":                              {"used_ledger"},
	"general_business_posting_group_id":        {"general_business_posting_group_id"},
	"general_product_posting_group_id":         {"general_product_posting_group_id"},
}
