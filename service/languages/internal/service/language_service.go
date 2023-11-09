package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type languageService struct {
	db                       *sql.DB
	simpleLanguageRepository treegrid.SimpleGridRowRepository
}

// GetPageCount implements languageService
func (o *languageService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return o.simpleLanguageRepository.GetPageCount(tr)
}

// GetPageData implements languageService
func (o *languageService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return o.simpleLanguageRepository.GetPageData(tr)
}

func NewLanguageService(db *sql.DB, simplelanguageService treegrid.SimpleGridRowRepository) LanguageService {
	return &languageService{db: db, simpleLanguageRepository: simplelanguageService}
}
