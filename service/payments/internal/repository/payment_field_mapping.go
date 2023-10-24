package repository

var (
	PaymentFieldNames = map[string][]string{
		"id":                        {"id"},
		"batch_id":                  {"batch_id"},
		"document_id":               {"document_id"},
		"document_no":               {"document_no"},
		"external_document_no":      {"external_document_no"},
		"transaction_no":            {"transaction_no"},
		"store_id":                  {"store_id"},
		"document_date":             {"document_date"},
		"posting_date":              {"posting_date"},
		"entry_date":                {"entry_date"},
		"account_type":              {"account_type"},
		"account_id":                {"account_id"},
		"recipient_bank_account_id": {"recipient_bank_account_id"},
		"balance_account_type":      {"balance_account_type"},
		"balance_account_id":        {"balance_account_id"},
		"amount":                    {"amount"},
		"amount_lcy":                {"amount_lcy"},
		"currency_id":               {"currency_id"},
		"currency_value":            {"currency_value"},
		"user_group_id":             {"user_group_id"},
		"status":                    {"status"},
		"payment_method_id":         {"payment_method_id"},
		"payment_reference":         {"payment_reference"},
		"creditor_no":               {"creditor_no"},
		"bank_payment_type_id":      {"bank_payment_type_id"},
		"bank_id":                   {"bank_id"},
		"paid":                      {"paid"},
		"remaining":                 {"remaining"},
		"paid_status":               {"paid_status"},
	}

	PaymentLineFieldNames = map[string][]string{
		"id":                    {"id"},
		"Parent":                {"parent_id"},
		"applies_document_type": {"applies_document_type"},
		"applies_document_id":   {"applies_document_id"},
		"payment_type_id":       {"payment_type_id"},
		"amount":                {"amount"},
		"amount_lcy":            {"amount_lcy"},
		"applied":               {"applied"},
	}

	//UserGroupLineFieldUploadNames = map[string][]string{
	//	"id":      {"user_group_lines.id"},
	//	"Parent":  {"user_group_lines.parent_id"}, // for easy parse
	//	"user_id": {"user_group_lines.user_id"},
	//}

	//UserUploadNames = map[string][]string{
	//	"email":     {"users.email"},
	//	"full_name": {"users.full_name"},
	//}
)
