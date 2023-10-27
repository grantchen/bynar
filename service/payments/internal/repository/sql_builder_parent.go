package repository

const (
	QueryParentCount = `
	SELECT COUNT(payments.id) as Count 
	FROM payments
	`

	QueryParent = `
	SELECT payments.id,
	payments.batch_id,
	payments.document_id,
	payments.document_no,
	payments.external_document_no,
	payments.transaction_no,
	payments.store_id,
	payments.document_date,
	payments.posting_date,
	payments.entry_date,
	payments.account_type,
	payments.account_id,
	payments.recipient_bank_account_id,
	payments.balance_account_type,
	payments.balance_account_id,
	payments.amount,
	payments.amount_lcy,
	payments.currency_id,
	payments.currency_value,
	payments.user_group_id,
	payments.payment_method_id,
	payments.payment_reference,
	payments.creditor_no,
	payments.bank_payment_type_id,
	payments.bank_id,
	payments.paid,
	payments.remaining,
	payments.paid_status,
	payments.status
	FROM payments
	`

	// empty
	QueryParentJoins = `
	`
)
