/**
    @author: dongjs
    @date: 2023/9/5
    @description:
**/

package models

type BaseURL string

type CheckoutScopes string

const (
	// todo move to checkout configuration
	GenerateAuthTokenURL BaseURL = `https://access.sandbox.checkout.com/connect/token`
	PaymentsURL          BaseURL = `https://api.sandbox.checkout.com/payments`

	GateWayPayment CheckoutScopes = "gateway:payment"
)
