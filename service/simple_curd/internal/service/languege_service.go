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
func (l *languageService) AddNewLanguage(newLang *model.Language) (*model.Language, error) {
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
func (l *languageService) DeleteLanguage(ID int) error {
	country, number, err := l.langRepository.GetCountryAndNumber(ID)
}

// GetAllLanguage implements LanguageService
func (l *languageService) GetAllLanguage() []*model.Language {
	panic("unimplemented")
}

// UpdateLanguage implements LanguageService
func (l *languageService) UpdateLanguage(lang *model.Language) (*model.Language, error) {
	panic("unimplemented")
}

func NewLanguageService(langRepository repository.LanguageRepository) LanguageService {
	return &languageService{langRepository: langRepository}
}
