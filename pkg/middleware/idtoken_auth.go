/**
    @author: dongjs
    @date: 2023/9/11
    @description:
**/

package middleware

import (
	"context"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// endpoints skip idToken auth
var skipIdTokenAuthEndEndpoints = []string{"/signin-email", "/signin", "/confirm-email", "/signup", "/verify-card", "/create-user"}

// http request header auth key
const httpAuthorizationHeader = "Authorization"

// get idToken from header of request
func getIdTokenFromHeader(r *http.Request) (error, string) {
	authorization := r.Header.Get(httpAuthorizationHeader)
	if "" != authorization && len(authorization) > 0 {
		tokens := strings.Split(authorization, " ")
		if len(tokens) == 2 {
			return nil, tokens[1]
		} else {
			return errors.New("token format error"), ""
		}
	}

	return errors.New("token is empty"), ""
}

// VerifyIdToken Verify idToken int http error code
func VerifyIdToken(r *http.Request) (int, string) {
	if utils.IsStringArrayInclude(skipIdTokenAuthEndEndpoints, r.RequestURI) {
		return http.StatusOK, ""
	}
	err, idToken := getIdTokenFromHeader(r)
	if err != nil {
		logrus.Errorf("get idToken from request error: %+v", err)
		return http.StatusBadRequest, ""
	}
	client, err := gip.NewGIPClient()
	if err != nil {
		logrus.Errorf("verifyIdToken: new GIPClient error: %v", err)
		return http.StatusInternalServerError, ""
	}
	_, err = client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		if err == gip.ErrIDTokenInvalid {
			return http.StatusUnauthorized, ""
		}
		return http.StatusInternalServerError, ""
	}
	//todo set current_user to request
	return http.StatusOK, ""
}
