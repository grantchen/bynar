package i18n

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"strings"
)

// Convert to 18 pieces of information based on err
func ErrMsgToI18n(err error, language string) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "of range") || strings.Contains(errMsg, "Truncated incorrect"):
		return fmt.Errorf(Localize(language, errors.ErrCodeOutRange))
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
