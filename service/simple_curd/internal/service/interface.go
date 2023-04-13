package service

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/model"

type LanguageService interface {
	GetAllLanguage() []*model.Language
	AddNewLanguage(newLang *model.Language) (*model.Language, error)
	UpdateLanguage(lang *model.Language) (*model.Language, error)
	DeleteLanguage(ID int) error
}
