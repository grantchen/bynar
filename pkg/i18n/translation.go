package i18n

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var errMsgToTranslationMap = map[string]string{
	"Truncated incorrect":         "TruncatedIncorrect",
	"INVALID_PHONE_NUMBER":        "PhoneNumberError",
	"phone number":                "PhoneNumberError",
	"phone_number already exists": "PhoneNumberAlreadyExists",
	"INVALID_EMAIL":               "EmailError",
	"email already exists":        "EmailAlreadyExists",
}

var globalBundle *i18n.Bundle

// No need to load active.en.toml since we are providing default translations.
// bundle.MustLoadMessageFile("active.en.toml")
func init() {
	globalBundle = i18n.NewBundle(language.English)
	globalBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// Read the configuration data in the i18n locales folder
	directoryPath := "../../pkg/i18n/tomls"
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		logrus.Errorf("i18n, error: %v", err)
	}
	for _, file := range files {
		fileName := file.Name()
		globalBundle.MustLoadMessageFile(directoryPath + "/" + fileName)
	}
}

func TranslationI18n(language, messageId string, templateData map[string]string) error {
	localizer := i18n.NewLocalizer(globalBundle, language)
	translationMessage := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    messageId,
			Other: "",
		},
		TemplateData: templateData,
	})

	e := NewError(language, translationMessage)
	e.withMessageId(messageId)
	return e
}

// TranslationErrorToI18n todo Database error translation
func TranslationErrorToI18n(language string, err error) error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}

	localizer := i18n.NewLocalizer(globalBundle, language)
	translationMessage := ""
	if err != nil {
		translationMessage = err.Error()
		for key, code := range errMsgToTranslationMap {
			if strings.Contains(translationMessage, key) {
				translationMessage = localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: code,
					},
				})
				break
			}
		}
	} else {
		return nil
	}

	e = NewError(language, translationMessage)
	e.WithCause(err)
	logrus.Errorf("%s, stack: %s", e.FullError(), e.Stack())
	return e
}
