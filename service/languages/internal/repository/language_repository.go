package repository

import "database/sql"

type languageRepository struct {
	db *sql.DB
}

func NewLanguageRepository(db *sql.DB) LanguageRepository {
	return &languageRepository{db: db}
}
