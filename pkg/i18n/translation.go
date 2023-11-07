package i18n

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var errMsgToTranslationMap = map[string]string{
	"Truncated incorrect":  "TruncatedIncorrect",
	"INVALID_PHONE_NUMBER": "PhoneNumberError",
	"phone number":         "PhoneNumberError",
	"INVALID_EMAIL":        "EmailError",
	"email":                "EmailError",
}

// No need to load active.en.toml since we are providing default translations.
// bundle.MustLoadMessageFile("active.en.toml")
func initBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// Read the configuration data in the i18n locales folder
	directoryPath := "../../pkg/i18n/tomls"
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		logrus.Errorf("i18n, error: %v", err)
	}
	for _, file := range files {
		fileName := file.Name()
		bundle.MustLoadMessageFile(directoryPath + "/" + fileName)
	}
	return bundle
}

func TranslationI18n(language, messageId string, templateData map[string]string) error {
	bundle := initBundle()
	localizer := i18n.NewLocalizer(bundle, language)
	translationMessage := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    messageId,
			Other: "",
		},
		TemplateData: templateData,
	})
	return fmt.Errorf(translationMessage)
}

// TranslationErrorToI18n todo Database error translation
func TranslationErrorToI18n(language string, err error) error {
	bundle := initBundle()
	localizer := i18n.NewLocalizer(bundle, language)
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
	return fmt.Errorf(translationMessage)
}
