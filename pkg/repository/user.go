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

// GetUserGroupID implements UserRepository
func (u *userRepository) GetUserGroupID(accountID int) (int, error) {
	var id int
	query := `
	SELECT parent_id 
	FROM user_group_lines
	WHERE user_id = ?
	`
	row := u.conn.QueryRow(query, accountID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// return new userRepository
func _(conn *sql.DB, accountID, moduleID int) UserRepository {
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

func (u *userRepository) GetUserID(_ int) (int, error) {
	return 0, nil
}
