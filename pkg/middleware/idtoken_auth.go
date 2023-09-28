/**
    @author: dongjs
    @date: 2023/9/11
    @description:
**/

package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"net/http"
	"strings"
)

// IdTokenClaims idToken decode struct
type IdTokenClaims struct {
	Name                string `json:"name"`
	OrganizationAccount bool   `json:"organization_account"`
	OrganizationStatus  bool   `json:"organization_status"`
	OrganizationUserId  int    `json:"organization_user_id"`
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
	TenantStatus        bool   `json:"tenant_status"`
	TenantSuspended     bool   `json:"tenant_suspended"`
	Firebase            struct {
		Identities struct {
			Email []string `json:"email"`
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
	Language  string `json:"language"`
	AccountId int    `json:"account_id"`
}

// endpoints skip idToken auth
var skipIdTokenAuthEndEndpoints = []string{"/signin-email", "/signin", "/confirm-email", "/signup", "/verify-card", "/create-user"}

// http request header auth key
const httpAuthorizationHeader = "Authorization"

type TokenAndDyDynamicDBContext struct {
	ConnectionString string
	DynamicDB        *sql.DB
	Claims           *IdTokenClaims
}

const IdTokenAndDynamicDBRequestContextKey string = "idTokenAndDynamicDB"

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
func VerifyIdToken(r *http.Request) (int, *IdTokenClaims, error) {
	if utils.IsStringArrayInclude(skipIdTokenAuthEndEndpoints, r.RequestURI) {
		return http.StatusOK, nil, nil
	}
	err, idToken := getIdTokenFromHeader(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	client, err := gip.NewGIPClient()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	claims, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		if err == gip.ErrIDTokenInvalid {
			return http.StatusUnauthorized, nil, err
		}
		return http.StatusInternalServerError, nil, err
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	var idTokenClaims = IdTokenClaims{}
	err = json.Unmarshal(claimsBytes, &idTokenClaims)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	// set current_user to request
	return http.StatusOK, &idTokenClaims, err
}

// GetIdTokenClaimsFromHttpRequestContext get idToken claims from request context
func GetIdTokenClaimsFromHttpRequestContext(r *http.Request) (*TokenAndDyDynamicDBContext, error) {
	tokenAndDyDynamicDBContext := r.Context().Value(IdTokenAndDynamicDBRequestContextKey)
	if tokenAndDyDynamicDBContext != nil {
		tokenAndDbContext := tokenAndDyDynamicDBContext.(*TokenAndDyDynamicDBContext)
		return tokenAndDbContext, nil
	}
	return nil, errors.New("get id_token fail")
}
