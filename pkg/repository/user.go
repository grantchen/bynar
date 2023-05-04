package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type userRepository struct {
	conn      *sql.DB
	accountID int
	moduleID  int
}

func NewUserRepository(conn *sql.DB, accountID, moduleID int) UserRepository {
	return &userRepository{conn: conn, accountID: accountID, moduleID: moduleID}
}

func (u *userRepository) AccountID() int {
	return u.accountID
}

func (u *userRepository) ModuleID() int {
	return u.moduleID
}

func (u *userRepository) HasPermission(moduleID string, accountID string) (bool, error) {
	var id int
	query := `
SELECT ugi.id, ugi.parent_id, ugi.account_id 
FROM user_group_items ugi
INNER JOIN user_groups ug ON ugi.parent_id = ug.id
INNER JOIN workflows w ON w.user_group_id = ug.id
WHERE w.module = ? AND ugi.account_id = ?
	`

	row := u.conn.QueryRow(query, moduleID, accountID)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("query row: [%w]", err)
	}

	return true, nil
}

func (u *userRepository) GetUserID(accountID int) (int, error) {
	return 0, nil
}
