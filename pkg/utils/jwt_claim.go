package utils

import "fmt"

type ErrMissingClaim struct {
	Key string
}

func (e ErrMissingClaim) Error() string {
	return fmt.Sprintf("missing claim '%s'", e.Key)
}

type ErrInvalidClaim struct {
	Key         string
	SourceError error
}

func (e ErrInvalidClaim) Error() string {
	return fmt.Sprintf("invalid claim '%s': %s", e.Key, e.SourceError)
}

func ResolveIntClaim(key string, claims map[string]interface{}) (int, error) {
	v, exists := claims[key]
	if !exists {
		return 0, &ErrMissingClaim{Key: key}
	}

	asInt, err := AsInt(v)
	if err != nil {
		return 0, &ErrInvalidClaim{Key: key, SourceError: err}
	}

	return asInt, err
}
