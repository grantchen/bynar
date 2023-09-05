/**
    @author: dongjs
    @date: 2023/9/5
    @description: checkout variables
**/

package configuration

import (
	"os"
)

type Environment interface {
	BaseUri() string
	AuthorizationUri() string
	IsSandbox() bool
}

type CheckoutEnv struct {
	baseUri          string
	authorizationUri string
	isSandbox        bool
}

func (e *CheckoutEnv) BaseUri() string {
	return e.baseUri
}

func (e *CheckoutEnv) AuthorizationUri() string {
	return e.authorizationUri
}

func (e *CheckoutEnv) IsSandbox() bool {
	return e.isSandbox
}

func NewEnvironment(
	baseUri string,
	authorizationUri string,
	isSandbox bool,
) *CheckoutEnv {
	return &CheckoutEnv{
		baseUri:          baseUri,
		authorizationUri: authorizationUri,
		isSandbox:        isSandbox}
}

// Sandbox test environment
func Sandbox() *CheckoutEnv {
	return NewEnvironment("https://api.sandbox.checkout.com",
		"https://access.sandbox.checkout.com/connect/token",
		true)
}

// Production product environment
func Production() *CheckoutEnv {
	return NewEnvironment(
		"https://api.checkout.com",
		"https://access.checkout.com/connect/token",
		false)
}

// CurrentEnv according .env file checkout.sandbox return CheckoutEnv
func CurrentEnv() *CheckoutEnv {
	sandbox := os.Getenv(ENVCheckoutSandBox)
	if "true" == sandbox || "" == sandbox {
		return Sandbox()
	} else {
		return Production()
	}

}
