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

type Environment interface {
	BaseUri() string
	AuthorizationUri() string
	PaymentsUri() string
	IsSandbox() bool
}

type CheckoutEnv struct {
	baseUri          string
	authorizationUri string
	paymentsUri      string
	isSandbox        bool
}

func (e *CheckoutEnv) BaseUri() string {
	return e.baseUri
}

func (e *CheckoutEnv) AuthorizationUri() string {
	return e.authorizationUri
}

func (e *CheckoutEnv) PaymentsUri() string {
	return e.paymentsUri
}

func (e *CheckoutEnv) IsSandbox() bool {
	return e.isSandbox
}

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
