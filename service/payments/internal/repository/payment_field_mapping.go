package repository

var (
	// PaymentFieldNames is a map for payments fields
	PaymentFieldNames = map[string][]string{
		"id":                        {"id", "payments.id"},
		"batch_id":                  {"batch_id", "payments.batch_id"},
		"document_id":               {"document_id", "payments.document_id"},
		"document_no":               {"document_no", "payments.document_no"},
		"external_document_no":      {"external_document_no", "payments.external_document_no"},
		"transaction_no":            {"transaction_no", "payments.transaction_no"},
		"document_date":             {"document_date", "payments.document_date"},
		"posting_date":              {"posting_date", "payments.posting_date"},
		"entry_date":                {"entry_date", "payments.entry_date"},
		"store_id":                  {"store_id", "payments.store_id"},
		"account_type":              {"account_type", "payments.account_type"},
		"account_id":                {"account_id", "payments.account_id"},
		"recipient_bank_account_id": {"recipient_bank_account_id", "payments.recipient_bank_account_id"},
		"balance_account_type":      {"balance_account_type", "payments.balance_account_type"},
		"balance_account_id":        {"balance_account_id", "payments.balance_account_id"},
		"currency_id":               {"currency_id", "payments.currency_id"},
		"currency_value":            {"currency_value", "payments.currency_value"},
		"amount":                    {"amount", "payments.amount"},
		"amount_lcy":                {"amount_lcy", "payments.amount_lcy"},
		"user_group_id":             {"user_group_id", "payments.user_group_id"},
		"payment_method_id":         {"payment_method_id", "payments.payment_method_id"},
		"payment_reference":         {"payment_reference", "payments.payment_reference"},
		"creditor_no":               {"creditor_no", "payments.creditor_no"},
		"bank_payment_type_id":      {"bank_payment_type_id", "payments.bank_payment_type_id"},
		"bank_id":                   {"bank_id", "payments.bank_id"},
		"status":                    {"status", "payments.status"},
		"paid":                      {"paid", "payments.paid"},
		"remaining":                 {"remaining", "payments.remaining"},
		"paid_status":               {"paid_status", "payments.paid_status"},
	}

	// PaymentLineFieldNames is a map for payment line fields
	PaymentLineFieldNames = map[string][]string{
		"id-line":               {"id", "payment_lines.id"},
		"Parent":                {"parent_id", "payment_lines.parent_id"},
		"applies_document_type": {"applies_document_type", "payment_lines.applies_document_type"},
		"applies_document_id":   {"applies_document_id", "payment_lines.applies_document_id"},
		"payment_type_id":       {"payment_type_id", "payment_lines.payment_type_id"},
		"amount":                {"amount", "payment_lines.amount"},
		"amount_lcy":            {"amount_lcy", "payment_lines.amount_lcy"},
		"applied":               {"applied", "payment_lines.applied"},
	}
	// PaymentFieldUploadNames is a map for user payment fields for upload
	PaymentFieldUploadNames = map[string][]string{
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
		"payment_method_id":         {"payment_method_id"},
		"payment_reference":         {"payment_reference"},
		"creditor_no":               {"creditor_no"},
		"bank_payment_type_id":      {"bank_payment_type_id"},
		"bank_id":                   {"bank_id"},
	}
	// PaymentLineFieldUploadNames is a map for user payment line fields for upload
	PaymentLineFieldUploadNames = map[string][]string{
		"id":                    {"id"},
		"Parent":                {"parent_id"},
		"applies_document_type": {"applies_document_type"},
		"applies_document_id":   {"applies_document_id"},
		"payment_type_id":       {"payment_type_id"},
		"amount":                {"amount"},
		"amount_lcy":            {"amount_lcy"},
		"applied":               {"applied"},
	}
)
