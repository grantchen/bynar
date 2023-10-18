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
	ErrCodeOutRange                 = "out-range"
	ErrCodePhoneNumber              = "phone-number"
	ErrCodeEmail                    = "email-error"
	ErrCodeTooLong                  = "too-long"
	ErrCodeGipUser                  = "gip-user-not-found"
	ErrCodeUserNotExist             = "user-not-exist"

	ErrCodeUserBelongSpecificUserGroupLines = "user-belong-specific-user-group-lines"
	ErrCodeNoUserGroupLineFound             = "no-user-group-line-found"
)

var (
	ErrCodeSave          = "save"
	ErrCodeUserGroup     = "user-group"
	ErrCodeUserGroupLine = "user-group-line"
)

var (
	ErrCodeArchivedUpdate           = "archived-update"
	ErrCodeArchivedDelete           = "archived-delete"
	ErrCodeArchivedNotValid         = "not-valid-archived"
	ErrCodeStatusNotValid           = "not-valid-status"
	ErrCodeSameArchivedStatus       = "status-and-archived-same"
	ErrCodeMergeRequest             = "merge-request"
	ErrCodeInvalidCondition         = "invalid-condition"
	ErrCodeNotField                 = "not-field-update"
	ErrCodeNoAllowToUpdateChildLine = "no-allow-to update-child-line"
)
