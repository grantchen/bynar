package i18n

import (
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
)

var errMsgToCodeMap = map[string]string{
	"duplicate":               errors.ErrCodeValueDuplicated,
	"gip user not found":      errors.ErrCodeGipUser,
	"user not found":          errors.ErrCodeNoUserFound,
	"not negative":            errors.ErrCodeNotNegativeNumber,
	"of range":                errors.ErrCodeOutRange,
	"Truncated incorrect":     errors.ErrCodeTruncatedIncorrect,
	"not field":               errors.ErrCodeNotField,
	"too long":                errors.ErrCodeTooLong,
	"Too Long":                errors.ErrCodeTooLong,
	"INVALID_PHONE_NUMBER":    errors.ErrCodePhoneNumber,
	"phone number":            errors.ErrCodePhoneNumber,
	"INVALID_EMAIL":           errors.ErrCodeEmail,
	"email":                   errors.ErrCodeEmail,
	"missing required params": errors.ErrCodeRequiredFieldsBlank,
}

// Deprecated: use TranslationI18n or TranslationErrorToI18n instead
func ErrMsgToI18n(err error, language string) error {
	errMsg := err.Error()
	for key, code := range errMsgToCodeMap {
		if strings.Contains(errMsg, key) {
			parts := strings.Split(errMsg, ",")
			message := Localize(language, code)
			if len(parts) > 1 {
				message += ": " + strings.Join(parts[1:], ": ")
			}
			return fmt.Errorf(message)
		}
	}
	return err
}
