package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var localesMap = map[string]map[string]string{}

// When you need to add a new language, just upload the configuration file for the new language in the locales folder
func init() {
	// Read the configuration data in the i18n locales folder
	locales, err := os.ReadDir("../../pkg/i18n/locales")
	if err != nil {
		logrus.Errorf("i18n, error: %v", err)
		return
	}
	for _, locale := range locales {
		fileName := locale.Name()
		localeByte, err := os.ReadFile("../../pkg/i18n/locales/" + fileName)
		if err != nil {
			logrus.Errorf("i18n, error: %v", err)
			return
		}

		localeMap := map[string]string{}
		if err = json.Unmarshal(localeByte, &localeMap); err != nil {
			logrus.Errorf("i18n, error: %v", err)
			return
		}

		// Get language name
		language := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
		localesMap[language] = localeMap
	}
}

// Localise the text based on the language passed from the front end
func Localize(r *http.Request, key string, args ...interface{}) string {
	language := r.Header.Get("x-language")

	// If the language is empty, the original text is returned
	// Compatible with scenarios that do not require i18n
	if language == "" {
		return key
	}

	// If it is an unsupported language, the empty text is returned
	localeMap, ok := localesMap[language]
	if !ok {
		logrus.Errorf("The language is not found: %s", language)
		return ""
	}

	// If no matching text is found, the empty text is returned
	msg, ok := localeMap[key]
	if !ok {
		logrus.Errorf("The key is not found: %s", key)
		return ""
	}

	if msg != "" {
		return fmt.Sprintf(msg, args...)
	}

	return msg
}
