/**
    @author: dongjs
    @date: 2023/9/6
    @description:
**/

package models

// AccessTokenResponse checkout generate access_token api response struct
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` //access token is valid for the length of time (in seconds) indicated by the expires_in field
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
