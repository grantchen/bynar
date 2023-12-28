package repository

const (
	QueryCount = `select count(*) from invoices inner join providers p on invoices.provider_id = p.id`

	QuerySelect = `select invoices.id, invoice_date, invoice_no, currency, total, billing_period_date, provider_id, paid from invoices inner join providers p on invoices.provider_id = p.id`

	AdditionWhere = `and account_id = %d`
)
