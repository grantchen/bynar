package errors

import "errors"

var (
	ErrForbiddenAction       = errors.New("forbidden action")
	ErrMissingRequiredParams = errors.New("missing required params")
	ErrInvalidQuantity       = errors.New("invalid quantity")
)
