/**
    @author: dongjs
    @date: 2023/9/5
    @description:
**/

package checkout

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/configuration"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"testing"
)

// generate_access_token test method
func TestGenerateAccessToken(t *testing.T) {
	// todo refactor
	client, err := checkout.NewPaymentClient()
	if err != nil {
		t.Log(err)
	}
	token, err := client.GenerateAuthToken(configuration.GatewayPayment)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("access_token is %+v", token)
	}
}

func TestValidateCard(t *testing.T) {
	// todo refactor
	client, err := checkout.NewPaymentClient()
	if err != nil {
		t.Log(err)
	}
	request := models.ValidateCardRequest{
		Token: "tok_tszdbkbr5bkulhgnuxho2gip3u",
		Email: "dongjs@tajansoft.com",
		Name:  "dongjinshuai",
	}
	card, err := client.ValidateCard(&request)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("card info is %+v", card)
	}
}
