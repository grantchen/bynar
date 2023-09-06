/**
    @author: dongjs
    @date: 2023/9/5
    @description: checkout payment api interface
**/

package checkout

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"

type PaymentClient interface {
	GenerateAuthToken(scope string) (models.AccessTokenResponse, error)
	ValidateCard(userDetails *models.ValidateCardRequest) (models.ValidateCard, error)
}
