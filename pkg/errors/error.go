package errors

import "errors"

var (
	ErrForbiddenAction       = errors.New("forbidden action")
	ErrMissingRequiredParams = errors.New("missing required params")
	ErrInvalidQuantity       = errors.New("invalid quantity")

	ErrSignInFail    = errors.New("sign in fail")
	ErrSendEmailFail = errors.New("send email fail")
	ErrNotSignUp     = errors.New("email not sign up")
)
