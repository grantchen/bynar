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

	"github.com/sirupsen/logrus"
)

// SendRegistrationEmail Send registration email
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
	resData, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		logrus.Error("Send email error: ", string(resData))
		return errors.New("failed to send email")
	}

	return nil
}

// SignInWithEmailLink Signs in a user with a out-of-band code from an email link.
func SignInWithEmailLink(email, oobCode string) error {
	oAuthClient, err := newOAuth2Client(context.Background())
	if err != nil {
		return err
	}

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithEmailLink?key=%s"
	url = fmt.Sprintf(url, os.Getenv(ENVGoogleAPIKey))
	data := map[string]interface{}{
		"oobCode": oobCode,
		"email":   email,
	}
	jsonByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	httpClient, err := oAuthClient.newHttpClient(context.Background())
	res, err := httpClient.Do(req)
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
