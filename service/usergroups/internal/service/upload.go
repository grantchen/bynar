package service

import "database/sql"

type UploadService struct {
	db *sql.DB
}

func NewUploadService(db *sql.DB) *UploadService {
	return &UploadService{
		db: db,
	}
}
