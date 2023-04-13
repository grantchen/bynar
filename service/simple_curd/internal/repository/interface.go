package repository

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/model"

type LanguageRepository interface {
	GetAllLanguage() []*model.Language
	AddNewLanguage(newLang *model.Language) (*model.Language, error)
	UpdateLanguage(lang *model.Language) (*model.Language, error)
	DeleteLanguage(ID int) error
	GetCountryAndNumber(id int) (country string, number int64, err error)
	ValidateOnIntegrity(id int, country string, number int64) (bool, error)
}
