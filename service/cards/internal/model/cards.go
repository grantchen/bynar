package model

type AccountsCard struct {
	ID                   int    `db:"id"`
	UserPaymentGatewayId string `db:"user_payment_gateway_id"`
	UserId               int    `db:"user_id"`
	Status               bool   `db:"status"`
	IsDefault            bool   `db:"is_default"`
	SourceId             int    `db:"source_id"`
	AccountId            int    `db:"account_id"`
}
