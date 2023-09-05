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
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/configuration"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type paymentClient struct {
	clientID            string
	clientSecret        string
	processingChannelID string
}

// GenerateAccessToken checkout payment generate access_token
func (p paymentClient) GenerateAccessToken(scope models.CheckoutScopes) (AccessTokenResponse, error) {
	var accessToken AccessTokenResponse
	payload := strings.NewReader("grant_type=client_credentials&scope=" + url.QueryEscape(string(scope)))
	client := &http.Client{}
	req, err := http.NewRequest("POST", configuration.CurrentEnv().AuthorizationUri(), payload)
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
func NewPaymentClient() (PaymentClient, error) {
	clientId := os.Getenv(configuration.ENVCheckoutClientId)
	if "" == clientId {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", configuration.ENVCheckoutClientId)
		err := errors.New(errMsg)
		log.Printf("NewPaymentClient: Error in getting environment variable %+v", err)
		return nil, err
	}
	clientSecret := os.Getenv(configuration.ENVCheckoutClientSecret)
	if "" == clientSecret {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", configuration.ENVCheckoutClientSecret)
		err := errors.New(errMsg)
		log.Printf("NewPaymentClient: Error in getting environment variable %+v", err)
		return nil, err
	}
	return &paymentClient{
		clientID:            clientId,
		clientSecret:        clientSecret,
		processingChannelID: "",
	}, nil
}
