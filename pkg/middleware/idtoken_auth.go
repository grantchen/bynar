/**
    @author: dongjs
    @date: 2023/9/11
    @description:
**/

package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// IdTokenClaims idToken decode struct
type IdTokenClaims struct {
	Name                string `json:"name"`
	OrganizationAccount bool   `json:"organization_account"`
	OrganizationStatus  bool   `json:"organization_status"`
	OrganizationUserId  string `json:"organization_user_id"`
	OrganizationUuid    string `json:"organization_uuid"`
	TenantUuid          string `json:"tenant_uuid"`
	Uid                 string `json:"uid"`
	Iss                 string `json:"iss"`
	Aud                 string `json:"aud"`
	AuthTime            int    `json:"auth_time"`
	UserId              string `json:"user_id"`
	Sub                 string `json:"sub"`
	Iat                 int    `json:"iat"`
	Exp                 int    `json:"exp"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"`
	Firebase            struct {
		Identities struct {
			Email []string `json:"email"`
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
}

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
func VerifyIdToken(r *http.Request) (int, string, *http.Request) {
	if utils.IsStringArrayInclude(skipIdTokenAuthEndEndpoints, r.RequestURI) {
		return http.StatusOK, "", r
	}
	err, idToken := getIdTokenFromHeader(r)
	if err != nil {
		logrus.Errorf("get idToken from request error: %+v", err)
		return http.StatusBadRequest, "", r
	}
	client, err := gip.NewGIPClient()
	if err != nil {
		logrus.Errorf("verifyIdToken: new GIPClient error: %v", err)
		return http.StatusInternalServerError, "", r
	}
	claims, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		logrus.Errorf("verifyIdToken: gip verify idToken error: %v", err)
		if err == gip.ErrIDTokenInvalid {
			return http.StatusUnauthorized, "", r
		}
		return http.StatusInternalServerError, "", r
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		logrus.Errorf("verifyIdToken: Marshal claims error: %v", err)
		return http.StatusInternalServerError, "", r
	}
	var idTokenClaims = IdTokenClaims{}
	err = json.Unmarshal(claimsBytes, &idTokenClaims)
	if err != nil {
		logrus.Errorf("verifyIdToken: Unmarshal claims error: %v", err)
		return http.StatusInternalServerError, "", r
	}
	// set current_user to request
	ctx := context.WithValue(context.Background(), "id_token", idTokenClaims)
	return http.StatusOK, "", r.WithContext(ctx)
}

// GetIdTokenClaimsFromHttpRequestContext get idToken claims from request context
func GetIdTokenClaimsFromHttpRequestContext(r *http.Request) (*IdTokenClaims, error) {
	idToken := r.Context().Value("id_token")
	if idToken != nil {
		claims := idToken.(IdTokenClaims)
		return &claims, nil
	}
	return nil, errors.New("no id_token fond in request context")
}
