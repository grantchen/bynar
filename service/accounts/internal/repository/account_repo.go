package repository

import (
	"github.com/sirupsen/logrus"
)

// CreateUser create a new account in db
func (r *accountRepositoryHandler) CreateUser(email string) error {
	_, err := r.db.Exec(`INSERT INTO accounts (email) VALUES (?);`, email)
	if err != nil {
		logrus.Errorf("CreateUser: error: %v", err)
		return err
	}
	return nil
}
