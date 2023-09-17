package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var locales = map[string]map[string]string{}

// When you need to add a new language, just upload the configuration file for the new language in the locales folder
func init() {
	// Read the configuration data in the i18n locales folder
	directoryPath := "../../pkg/i18n/locales"
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		logrus.Errorf("i18n, error: %v", err)
		return
	}
	for _, file := range files {
		fileName := file.Name()
		fileByte, err := os.ReadFile(directoryPath + "/" + fileName)
		if err != nil {
			logrus.Errorf("i18n, error: %v", err)
			return
		}

		locale := map[string]string{}
		if err = json.Unmarshal(fileByte, &locale); err != nil {
			logrus.Errorf("i18n, error: %v", err)
			return
		}

		// Get language name
		language := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
		locales[language] = locale
	}
}

// Localise the text based on the language
func Localize(language, key string, args ...interface{}) string {
	// If the language is empty, English is used by default
	if language == "" {
		language = "en"
	}

	// If it is an unsupported language, the empty text is returned
	locale, ok := locales[language]
	if !ok {
		logrus.Errorf("The language is not found: %s", language)
		return ""
	}

	// If no matching text is found, the empty text is returned
	msg, ok := locale[key]
	if !ok {
		logrus.Errorf("The key is not found: %s", key)
		return ""
	}

	if msg != "" {
		return fmt.Sprintf(msg, args...)
	}

	return msg
}
