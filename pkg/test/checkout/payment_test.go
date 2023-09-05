/**
    @author: dongjs
    @date: 2023/9/5
    @description:
**/

package checkout

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"github.com/joho/godotenv"
	"log"
	"testing"
)

// generate_access_token test method
func TestGenerateAccessToken(t *testing.T) {
	// todo test
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	client, err := checkout.NewPaymentClient()
	if err != nil {
		t.Log(err)
	}
	token, err := client.GenerateAccessToken(models.GateWayPayment)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("access_token is %+v", token)
	}

}
