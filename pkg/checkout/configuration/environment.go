/**
    @author: dongjs
    @date: 2023/9/5
    @description: checkout variables
**/

package configuration

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/constant"
	"os"
)

// Environment checkout env interface
type Environment interface {
	// BaseUri checkout base url
	BaseUri() string
	// AuthorizationUri checkout connect token url
	AuthorizationUri() string
	// PaymentsUri checkout payments url
	PaymentsUri() string
	// IsSandbox checkout sandbox or not
	IsSandbox() bool
}

// CheckoutEnv checkout env config
type CheckoutEnv struct {
	baseUri          string
	authorizationUri string
	paymentsUri      string
	isSandbox        bool
}

// BaseUri return checkout base uri
func (e *CheckoutEnv) BaseUri() string {
	return e.baseUri
}

// AuthorizationUri return checkout authorization api uri
func (e *CheckoutEnv) AuthorizationUri() string {
	return e.authorizationUri
}

// PaymentsUri return checkout payments api uri
func (e *CheckoutEnv) PaymentsUri() string {
	return e.paymentsUri
}

// IsSandbox return current env is sandbox test
func (e *CheckoutEnv) IsSandbox() bool {
	return e.isSandbox
}

// NewEnvironment create checkout env
func NewEnvironment(
	baseUri string,
	authorizationUri string,
	paymentsUri string,
	isSandbox bool,
) *CheckoutEnv {
	return &CheckoutEnv{
		baseUri:          baseUri,
		authorizationUri: authorizationUri,
		paymentsUri:      paymentsUri,
		isSandbox:        isSandbox}
}

// Sandbox test environment
func Sandbox() *CheckoutEnv {
	return NewEnvironment("https://api.sandbox.checkout.com",
		"https://access.sandbox.checkout.com/connect/token",
		"https://api.sandbox.checkout.com/payments",
		true)
}

// Production product environment
func Production() *CheckoutEnv {
	return NewEnvironment(
		"https://api.checkout.com",
		"https://access.checkout.com/connect/token",
		"https://api.checkout.com/payments",
		false)
}

// CurrentEnv according .env file checkout.sandbox return CheckoutEnv
func CurrentEnv() *CheckoutEnv {
	sandbox := os.Getenv(constant.ENVCheckoutSandBox)
	if "true" == sandbox || "" == sandbox {
		return Sandbox()
	} else {
		return Production()
	}

}