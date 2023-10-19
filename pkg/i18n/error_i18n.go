package i18n

import (
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
)

// Convert to 18 pieces of information based on err
func ErrMsgToI18n(err error, language string) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "missing required params"):
		return fmt.Errorf(strings.ReplaceAll(errMsg, "missing required params", Localize(language, "missing-required-params")))
	case strings.Contains(errMsg, "of range"):
		return fmt.Errorf(Localize(language, errors.ErrCodeOutRange))
	case strings.Contains(errMsg, "Truncated incorrect"):
		return fmt.Errorf(Localize(language, errors.ErrCodeTruncatedIncorrect))
	case strings.Contains(errMsg, "not field"):
		return fmt.Errorf(Localize(language, errors.ErrCodeNotField))
	case strings.Contains(errMsg, "too long") || strings.Contains(errMsg, "Too Long"):
		return fmt.Errorf(Localize(language, errors.ErrCodeTooLong))
	case strings.Contains(errMsg, "INVALID_PHONE_NUMBER") || strings.Contains(errMsg, "phone number"):
		return fmt.Errorf(Localize(language, errors.ErrCodePhoneNumber))
	case strings.Contains(errMsg, "INVALID_EMAIL") || strings.Contains(errMsg, "email"):
		return fmt.Errorf(Localize(language, errors.ErrCodeEmail))
	default:
		return err
	}
}
