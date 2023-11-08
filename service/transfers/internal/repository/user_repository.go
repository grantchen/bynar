package repository

import (
	"database/sql"
)

// UserRepository is a repository for user
type userRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{}
}
