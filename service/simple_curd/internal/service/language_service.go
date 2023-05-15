package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/repository"
)

type languageService struct {
	langRepository repository.LanguageRepository
}

var (
	ErrMissingRequiredParams = errors.New("missing required params")
	ErrAlreadyExist          = errors.New("already exist")
)

// AddNewLanguage implements LanguageService
func (l *languageService) AddNewLanguage(row *model.Changes) (*model.Language, error) {
	newLang, err := gridRowToLanguage(row)
	if err != nil {
		return nil, fmt.Errorf("Parse param error: [%w]", err)
	}

	ok, err := l.langRepository.ValidateOnIntegrity(newLang.Id, newLang.Country, newLang.Number)
	if err != nil {
		return nil, fmt.Errorf("validate on integrity: [%w]", err)

	}

	if !ok {
		return nil, fmt.Errorf("[%w]: country '%s', code '%d'", ErrAlreadyExist, newLang.Country, newLang.Number)
	}

	lang, err := l.langRepository.AddNewLanguage(newLang)
	return lang, err
}

func gridRowToLanguage(row *model.Changes) (*model.Language, error) {
	id, err := strconv.Atoi(row.Id)
	lang := &model.Language{
		Country:       row.Country,
		Language:      row.Language,
		Two_letters:   row.TwoLetters,
		Three_letters: row.ThreeLetters,
	}
	if err == nil {
		lang.Id = id
	}

	number, err := strconv.ParseInt(row.Number, 10, 64)
	if err == nil {
		lang.Number = number
	}

	return lang, nil
}

func validateOnRequiredAll(row *model.Changes) error {
	switch {
	case strings.TrimSpace(row.Country) == "":
		return fmt.Errorf("[%w]: %s", ErrMissingRequiredParams, "country")
	case strings.TrimSpace(row.Number) == "":
		return fmt.Errorf("[%w]: %s", ErrMissingRequiredParams, "number")
	case strings.TrimSpace(row.Language) == "":
		return fmt.Errorf("[%w]: %s", ErrMissingRequiredParams, "language")
	case strings.TrimSpace(row.TwoLetters) == "":
		return fmt.Errorf("[%w]: %s", ErrMissingRequiredParams, "two_letters")
	case strings.TrimSpace(row.ThreeLetters) == "":
		return fmt.Errorf("[%w]: %s", ErrMissingRequiredParams, "three_letters")
	}

	return nil
}

// DeleteLanguage implements LanguageService
func (l *languageService) DeleteLanguage(ID int64) error {
	return l.langRepository.DeleteLanguage(ID)
}

// GetAllLanguage implements LanguageService
func (l *languageService) GetAllLanguage() []*model.Language {
	return l.langRepository.GetAllLanguage()
}

// UpdateLanguage implements LanguageService
func (l *languageService) UpdateLanguage(row *model.Changes) (*model.Language, error) {
	id, _ := strconv.Atoi(row.Id)
	country, number, err := l.langRepository.GetCountryAndNumber(id)
	if err != nil {
		return nil, fmt.Errorf("get country and number (id '%s'): [%w]", row.Id, err)
	}

	// if not empty then should update 'country' for checking on integrity
	if row.Country != "" {
		country = row.Country
	}

	// if not empty then should update 'number' for checking on integrity
	if row.Number != "" {
		number, _ = strconv.ParseInt(row.Number, 10, 64)
	}

	// check if another row with the same 'number' and 'country' exist
	// ok = false - not valid means validation didn't pass and row exist
	ok, err := l.langRepository.ValidateOnIntegrity(id, country, number)
	if err != nil {
		return nil, fmt.Errorf("validate on integrity: [%w]", err)
	}

	// ok = false - validation failed, return ErrAlreadyExist
	if !ok {
		return nil, fmt.Errorf("[%w]: id '%s', country '%s', number '%d'",
			ErrAlreadyExist, row.Id, country, number)
	}

	lang, err := gridRowToLanguage(row)
	if err != nil {
		return nil, fmt.Errorf("parse param error: [%w]", err)
	}

	return l.langRepository.UpdateLanguage(lang)

}

func NewLanguageService(langRepository repository.LanguageRepository) LanguageService {
	return &languageService{langRepository: langRepository}
}
