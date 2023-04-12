package mapping

var (
	CashReceiptFieldNames = map[string][]string{
		"id":                   {"id"},
		"batch_id":             {"batch_id"},
		"document_id":          {"document_id"},
		"document_no":          {"document_no"},
		"transaction_no":       {"transaction_no"},
		"store_id":             {"store_id"},
		"document_date":        {"document_date"},
		"posting_date":         {"posting_date"},
		"entry_date":           {"entry_date"},
		"account_type":         {"account_type"},
		"account_id":           {"account_id"},
		"balance_account_type": {"balance_account_type"},
		"balance_account_id":   {"balance_account_id"},
		"amount":               {"amount"},
		"amount_lcy":           {"amount_lcy"},
		"currency_value":       {"currency_value"},
		"user_group_id":        {"user_group_id"},
		"status":               {"status"},
		"bank_id":              {"bank_id"},
	}

	CashReceiptLineFieldNames = map[string][]string{
		"id":                    {"id"},
		"Parent":                {"parent_id"},
		"applies_document_type": {"applies_document_type"},
		"applies_document_id":   {"applies_document_id"},
		"amount":                {"amount"},
		"amount_lcy":            {"amount_lcy"},
		"applied":               {"applied"},
	}
)
