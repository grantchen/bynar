package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidPayload     = errors.New("invalid payload")
)

const (
	BearerTokenType = "Bearer"

	tokenSegmentSeparator = "."
)

func ExtractClaimsFromJWT(token string) (claims map[string]interface{}, _ error) {
	payload, err := extractPayload(token)
	if err != nil {
		return nil, err
	}

	if unmarshalErr := json.Unmarshal(payload, &claims); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return claims, nil
}

func extractPayload(token string) ([]byte, error) {
	token = strings.TrimPrefix(token, BearerTokenType+" ")
	token = strings.TrimPrefix(token, strings.ToLower(BearerTokenType)+" ")

	parts := strings.Split(token, tokenSegmentSeparator)
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidPayload
	}

	return payload, nil
}

func UnmarshalClaimsFromJWT(token string, output interface{}) error {
	payload, err := extractPayload(token)
	if err != nil {
		return err
	}

	return json.Unmarshal(payload, output)
}
