package service

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/model"

type LanguageService interface {
	GetAllLanguage() []*model.Language
	AddNewLanguage(row *model.Changes) (*model.Language, error)
	UpdateLanguage(row *model.Changes) (*model.Language, error)
	DeleteLanguage(ID int64) error
}
