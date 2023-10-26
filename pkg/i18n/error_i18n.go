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

func ErrMsgToI18n(err error, language string) error {
	errMsg := err.Error()
	for key, code := range errMsgToCodeMap {
		if strings.Contains(errMsg, key) {
			parts := strings.SplitN(errMsg, ",", 2)
			message := Localize(language, code)
			if len(parts) > 1 {
				message += ": " + strings.TrimSpace(parts[0])
			}
			return fmt.Errorf(message)
		}
	}
	return err
}
