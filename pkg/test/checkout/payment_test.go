/**
    @author: dongjs
    @date: 2023/9/5
    @description:
**/

package checkout

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/configuration"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
)

var client checkout.PaymentClient

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../service/main/.env")
	if err != nil {
		log.Fatal(err)
	}
	newClient, err := checkout.NewPaymentClient()
	if err != nil {
		log.Fatalln(err)
	}
	client = newClient
	os.Exit(m.Run())
}

// generate_access_token test method
func TestGenerateAccessToken(t *testing.T) {
	token, err := client.GenerateAuthToken(configuration.GatewayPayment)
	if err != nil {
		t.Fatalf("TestGenerateAccessToken: generate auth token error is %+v", err)
	} else {
		t.Logf("access_token is %+v", token)
	}
}

func TestValidateCard(t *testing.T) {
	request := models.ValidateCardRequest{
		Token: "tok_tszdbkbr5bkulhgnuxho2gip3u",
		Email: "dongjs@tajansoft.com",
		Name:  "dongjinshuai",
	}
	card, err := client.ValidateCard(&request)
	if err != nil {
		t.Fatalf("TestValidateCard: validte card error is %+v", err)
	} else {
		t.Logf("card info is %+v", card)
	}
}

func TestFetchCustomerDetails(t *testing.T) {
	res, err := client.FetchCustomerDetails("cus_4xeunchcg3lu3gcyh222qfgvr4")
	if err != nil {
		t.Fatalf("TestFetchCustomerDetails: fetch customer details error is %+v", err)
	} else {
		t.Logf("customer detail is %+v", res)
	}

}
