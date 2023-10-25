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
	case strings.Contains(errMsg, "duplicate"):
		parts := strings.Split(errMsg, ",")
		return fmt.Errorf("%s: %s", parts[0], Localize(language, errors.ErrCodeValueDuplicated))
	case strings.Contains(errMsg, "gip user not found"):
		return fmt.Errorf(Localize(language, errors.ErrCodeGipUser))
	case strings.Contains(errMsg, "user not found"):
		return fmt.Errorf(Localize(language, errors.ErrCodeNoUserFound))
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
	case strings.Contains(errMsg, "missing required params"):
		index := strings.Index(errMsg, ":")
		result := ""
		if index != -1 && index+1 < len(errMsg) {
			result = err.Error()[index+1:]
		}
		return fmt.Errorf("%s: %s", Localize(language, errors.ErrCodeRequiredFieldsBlank), result)
	default:
		return err
	}
}
