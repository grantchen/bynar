/**
    @author: dongjs
    @date: 2023/9/5
    @description: checkout api
**/

package checkout

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type paymentClient struct {
	sm                  secretsmanager.SecretsManager // todo remove
	clientID            string
	clientSecret        string
	processingChannelID string
}

// GenerateAccessToken checkout payment generate access_token
func (p paymentClient) GenerateAccessToken(scope models.CheckoutScopes) (AccessTokenResponse, error) {
	var accessToken AccessTokenResponse
	payload := strings.NewReader("grant_type=client_credentials&scope=" + url.QueryEscape(string(scope)))
	client := &http.Client{}
	req, err := http.NewRequest("POST", string(models.GenerateAuthTokenURL), payload)
	if err != nil {
		log.Printf("GenerateAuthToken: Error in creating new request %v", err)
		return accessToken, err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(p.clientID + ":" + p.clientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+encoded)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("GenerateAuthToken: Error in generating access token %+v", err)
		return accessToken, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		var errResp interface{}
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			log.Printf("GenerateAuthToken: Error in decoding error response %+v", err)
			return accessToken, err
		}
		log.Printf("GenerateAuthToken: Error in validating card %+v", errResp)
		return accessToken, errors.New("error in generating token")
	}
	err = json.NewDecoder(res.Body).Decode(&accessToken)
	if err != nil {
		log.Printf("GenerateAuthToken: Error in decoding response %+v", err)
		return accessToken, err
	}
	return accessToken, nil
}

// NewPaymentClient create paymentClient to call checkout api
func NewPaymentClient(sm secretsmanager.SecretsManager) (PaymentClient, error) {
	// todo test set sandbox clientId and clientSecret and clientId etc will set to .env file and will remove input params sm
	return &paymentClient{
		sm:                  sm,
		clientID:            "ack_3kgxgdj773yubf4sfmiht3r4h4",
		clientSecret:        "PddTMk1FBjk1MDQHtBt1U8cHjZvS+Guc80NmcUHp3pHevOpt7EgYkT/DWae7gnOTlF6kPCPo+RZEu9xut/5VWA==",
		processingChannelID: "",
	}, nil
}
