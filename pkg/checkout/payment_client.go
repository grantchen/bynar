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

// checkout payment client struct
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
		return nil, errors.New(errMsg)
	}
	clientSecret := os.Getenv(constant.ENVCheckoutClientSecret)
	if "" == clientSecret {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", constant.ENVCheckoutClientSecret)
		return nil, errors.New(errMsg)
	}
	processingChannelID := os.Getenv(constant.ENVCheckoutProcessChannelId)
	if "" == processingChannelID {
		var errMsg = fmt.Sprintf("no %s variable in .env file or blank", constant.ENVCheckoutProcessChannelId)
		err := errors.New(errMsg)
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
		return accessToken, err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(p.clientID + ":" + p.clientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+encoded)

	res, err := client.Do(req)
	if err != nil {
		return accessToken, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		var errResp models.AccessTokenErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return accessToken, err
		}
		log.Printf("GenerateAuthToken: Error in validating card %+v", errResp)
		return accessToken, errors.New(errResp.Error)
	}
	err = json.NewDecoder(res.Body).Decode(&accessToken)
	if err != nil {
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
		return resp, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, configuration.CurrentEnv().PaymentsUri(), bytes.NewReader(payloadBytes))
	if err != nil {
		return resp, err
	}

	authorization, err := p.GenerateAuthToken(configuration.GatewayPayment)
	if err != nil {
		return resp, err
	}
	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 && res.StatusCode != 202 {
		var errResp models.CheckOutErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return resp, err
		}
		logrus.Errorf("ValidateAndStoreCard: Error in validating card %+v", errResp)
		return resp, errors.New(strings.Join(errResp.ErrorCodes, ";"))
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DeleteCard Delete an instrument
func (p paymentClient) DeleteCard(sourceID string) error {
	apiURL := fmt.Sprintf(`%v/%v`, configuration.CurrentEnv().InstrumentUri(), sourceID)
	method := "DELETE"
	authorization, err := p.GenerateAuthToken(configuration.VaultInstruments)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, nil)
	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		var errResp models.CheckOutErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return err
		}
		logrus.Errorf("DeleteCard: Error in delete card %+v", errResp)
		return errors.New(strings.Join(errResp.ErrorCodes, ";"))
	}
	return nil
}

// DeleteCustomer Delete a customer and all of their linked payment instruments.
func (p paymentClient) DeleteCustomer(customerID string) error {
	apiURL := fmt.Sprintf(`%v/%v`, configuration.CurrentEnv().CustomerUri(), customerID)
	method := "DELETE"
	authorization, err := p.GenerateAuthToken(configuration.Vault)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, nil)
	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		var errResp models.CheckOutErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return err
		}
		logrus.Errorf("DeleteCustomer: Error in delete customer %+v", errResp)
		return errors.New(strings.Join(errResp.ErrorCodes, ";"))
	}
	return nil
}

// UpdateCustomer updates customer information in checkout like name, email and default card
func (p paymentClient) UpdateCustomer(customerInfo models.UpdateCustomer, customerID string) error {
	apiURL := fmt.Sprintf(`%v/%v`, configuration.CurrentEnv().CustomerUri(), customerID)
	method := "PATCH"
	payload := strings.NewReader(fmt.Sprintf(`{
						"email": "%v",
						"name": "%v",
						"default": "%v"
					  }`, customerInfo.Email, customerInfo.Name, customerInfo.DefaultInstrument))

	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, payload)
	if err != nil {
		logrus.Errorf("UpdateCustomer: Error creating new request %v", err)
		return err
	}
	authorization, err := p.GenerateAuthToken(configuration.VaultInstruments)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		var errResp models.CheckOutErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return err
		}
		logrus.Errorf("DeleteCard: Error in delete card %+v", errResp)
		return errors.New(strings.Join(errResp.ErrorCodes, ";"))
	}

	return nil
}

// FetchPaymentDetails fetch payment details of customer
func (p paymentClient) FetchPaymentDetails(paymentID string) (paymentDetails models.FetchPaymentDetails, err error) {
	apiURL := fmt.Sprintf(`%v/%v`, configuration.CurrentEnv().PaymentsUri(), paymentID)
	method := "GET"
	authorization, err := p.GenerateAuthToken(configuration.GatewayPaymentDetails)
	if err != nil {
		return paymentDetails, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, nil)
	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return paymentDetails, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errResp models.CheckOutErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return paymentDetails, err
		}
		logrus.Errorf("GetCustomerDetails: Error in fetching payment details %+v with statusCode %v", errResp, res.StatusCode)
		return paymentDetails, errors.New(strings.Join(errResp.ErrorCodes, ";"))
	}

	err = json.NewDecoder(res.Body).Decode(&paymentDetails)
	if err != nil {
		return paymentDetails, err
	}

	return paymentDetails, nil
}

// FetchCustomerDetails fetch customer all card details
func (p paymentClient) FetchCustomerDetails(customerID string) (models.CustomerResponse, error) {
	var customerDetails models.CustomerResponse
	apiURL := fmt.Sprintf(`%v/%v`, configuration.CurrentEnv().CustomerUri(), customerID)
	method := "GET"
	authorization, err := p.GenerateAuthToken(configuration.VaultInstruments)
	if err != nil {
		return customerDetails, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, nil)
	req.Header.Add("Authorization", "Bearer "+authorization.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return customerDetails, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errResp models.CheckOutErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return customerDetails, err
		}
		logrus.Errorf("GetCustomerDetails: Error in fetching customer details %+v with statusCode %v", errResp, res.StatusCode)
		return customerDetails, errors.New(strings.Join(errResp.ErrorCodes, ";"))
	}

	err = json.NewDecoder(res.Body).Decode(&customerDetails)
	if err != nil {
		return customerDetails, err
	}

	return customerDetails, nil
}
