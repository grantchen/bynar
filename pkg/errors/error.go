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

var (
	ErrCode                         = "error"
	ErrCodePermissionDenied         = "permission-denied"
	ErrCodeMethodNotAllowed         = "method-no-allowed"
	ErrCodeUploadProfileFail        = "upload-profile-fail"
	ErrCodeDeleteProfileFail        = "delete-profile-fail"
	ErrCodeOperationFailed          = "operation-failed"
	ErrCodeRequestParameter         = "request-parameter-error"
	ErrCodeNoUserFound              = "no-user-found"
	ErrCodeEmailAlreadyExists       = "email-already-exists"
	ErrCodePhoneNumberAlreadyExists = "phone_number-already-exists"
	ErrCodeIDTokenInvalid           = "id-token-invalid"
	ErrCodeSignatureInvalid         = "signature-invalid"
	ErrCodeTimestampExpired         = "timestamp-expired"
	ErrCodeValueDuplicated          = "value-duplicate"
	ErrCodeRequiredFieldsBlank      = "required-fields-blank"
	ErrCodeUndefinedTowType         = "undefined-row-type"
	ErrCodeBeginTransaction         = "begin-transaction"
	ErrCodeCommitTransaction        = "commit-transaction"
)
