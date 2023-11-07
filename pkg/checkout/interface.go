/**
    @author: dongjs
    @date: 2023/9/5
    @description: checkout payment api interface
**/

package checkout

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"

// PaymentClient checkout api interface
type PaymentClient interface {
	// GenerateAuthToken generate payment api access token
	GenerateAuthToken(scope string) (models.AccessTokenResponse, error)
	// ValidateCard validate card api
	ValidateCard(userDetails *models.ValidateCardRequest) (models.ValidateCard, error)
	FetchCustomerDetails(customerID string) (models.CustomerResponse, error)
	DeleteCard(sourceID string) error
	// DeleteCustomer Delete a customer and all of their linked payment instruments.
	DeleteCustomer(customerID string) error
	UpdateCustomer(customerInfo models.UpdateCustomer, customerID string) error
	FetchPaymentDetails(paymentID string) (models.FetchPaymentDetails, error)
}
