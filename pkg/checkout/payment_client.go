/**
    @author: dongjs
    @date: 2023/9/5
    @description: checkout api
**/

package checkout

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/configuration"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/constant"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"github.com/sirupsen/logrus"
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

// NewPaymentClient create paymentClient to call checkout api
func NewPaymentClient() (PaymentClient, error) {
	clientId := os.Getenv(constant.ENVCheckoutClientId)
	if "" == clientId {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", constant.ENVCheckoutClientId)
		err := errors.New(errMsg)
		log.Printf("NewPaymentClient: Error in getting environment variable %+v", err)
		return nil, err
	}
	clientSecret := os.Getenv(constant.ENVCheckoutClientSecret)
	if "" == clientSecret {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", constant.ENVCheckoutClientSecret)
		err := errors.New(errMsg)
		log.Printf("NewPaymentClient: Error in getting environment variable %+v", err)
		return nil, err
	}
	processingChannelID := os.Getenv(constant.ENVCheckoutProcessChannelId)
	if "" == processingChannelID {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", constant.ENVCheckoutProcessChannelId)
		err := errors.New(errMsg)
		log.Printf("NewPaymentClient: Error in getting environment variable %+v", err)
		return nil, err
	}
	return &paymentClient{
		clientID:            clientId,
		clientSecret:        clientSecret,
		processingChannelID: processingChannelID,
	}, nil
}

// GenerateAuthToken checkout payment generate access_token
func (p paymentClient) GenerateAuthToken(scope string) (models.AccessTokenResponse, error) {
	var accessToken models.AccessTokenResponse
	payload := strings.NewReader("grant_type=client_credentials&scope=" + url.QueryEscape(scope))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, configuration.CurrentEnv().AuthorizationUri(), payload)
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

// ValidateCard validate user card details
func (p paymentClient) ValidateCard(userDetails *models.ValidateCardRequest) (models.ValidateCard, error) {
	var resp models.ValidateCard
	payload := models.CardValidationPayload{
		Source: models.TokenSource{
			Type:  "token",
			Token: userDetails.Token,
		},
		Currency: "USD",
		Customer: models.NewCustomer{
			Email: userDetails.Email,
			Name:  userDetails.Name,
		},
		ProcessingChannelID: p.processingChannelID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("ValidateAndStoreCard: Error in marshaling payload %+v", err)
		return resp, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, configuration.CurrentEnv().PaymentsUri(), bytes.NewReader(payloadBytes))
	if err != nil {
		logrus.Errorf("ValidateAndStoreCard: Error in creating new request %v", err)
		return resp, err
	}

	authorization, err := p.GenerateAuthToken(configuration.GatewayPayment)
	if err != nil {
		logrus.Errorf("ValidateAndStoreCard: Error in generating auth token %v", err)
		return resp, err
	}
	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logrus.Errorf("ValidateAndStoreCard: Error in creating 0$ payment %+v", err)
		return resp, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 && res.StatusCode != 202 {
		var errResp interface{}
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			logrus.Errorf("ValidateAndStoreCard: Error in decoding error response %+v", err)
			return resp, err
		}
		logrus.Errorf("ValidateAndStoreCard: Error in validating card %+v", errResp)
		return resp, errors.New("error in validating card")
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.Errorf("ValidateAndStoreCard: Error in decoding response %+v", err)
		return resp, err
	}

	return resp, nil
}
