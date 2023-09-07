package gip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Send registration email
func SendRegistrationEmail(email, continueUrl string) error {
	oAuthClient, err := newOAuth2Client(context.Background())
	if err != nil {
		return err
	}

	url := "https://identitytoolkit.googleapis.com/v1/projects/%s/accounts:sendOobCode?key=%s"
	url = fmt.Sprintf(url, oAuthClient.projectID, os.Getenv(ENVGoogleAPIKey))
	data := map[string]interface{}{
		"requestType": "EMAIL_SIGNIN",
		"email":       email,
		"continueUrl": continueUrl,
	}
	jsonByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	httpClient, err := oAuthClient.newHttpClient(context.Background())
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, _ = io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return errors.New("failed to send email")
	}

	return nil
}

// Verification email, invalid for type 'EMAIL SIGNIN'
func VerificationEmail(oobCode string) error {
	url := "https://identitytoolkit.googleapis.com/v1/accounts:update?key=%s"
	url = fmt.Sprintf(url, os.Getenv(ENVGoogleAPIKey))
	data := map[string]interface{}{
		"oobCode": oobCode,
	}
	jsonByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return errors.New(string(response))
	}

	return nil
}
