package repository

const (
	QueryCount = `SELECT COUNT(*) as Count FROM general_posting_setup`

	QuerySelect = `SELECT * FROM general_posting_setup
	INNER JOIN general_business_posting_groups gbpg on general_business_posting_group_id = gbpg.id
	INNER JOIN general_product_posting_groups gppg on general_posting_setup.general_product_posting_group_id = gppg.id`

	QueryJoin = `INNER JOIN general_business_posting_groups gbpg on general_business_posting_group_id = gbpg.id
	INNER JOIN general_product_posting_groups gppg on general_posting_setup.general_product_posting_group_id = gppg.id`

	QuerySelectWithoutJoin = `SELECT * FROM general_posting_setup`
)
