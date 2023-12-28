package errors

import "errors"

var (
	ErrForbiddenAction = errors.New("forbidden action")
)

var (
	ErrCode                         = "error"
	ErrCodeUploadProfileFail        = "upload-profile-fail"
	ErrCodeNoUserFound              = "no-user-found"
	ErrCodeEmailAlreadyExists       = "email-already-exists"
	ErrCodePhoneNumberAlreadyExists = "phone_number-already-exists"
	ErrCodeIDTokenInvalid           = "id-token-invalid"
	ErrCodeSignatureInvalid         = "signature-invalid"
	ErrCodeTimestampExpired         = "timestamp-expired"
)
