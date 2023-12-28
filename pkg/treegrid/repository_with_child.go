package treegrid

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

type GridRowRepositoryWithChild interface {
	IsChild(gr GridRow) bool
	GetParentID(gr GridRow) (parentID interface{}, err error)
	Save(tx *sql.Tx, tr *MainRow) error
	SaveMainRow(tx *sql.Tx, tr *MainRow) error
	SaveLines(tx *sql.Tx, tr *MainRow) error

	// SaveLineAdd use for custom save lines.
	SaveLineAdd(tx *sql.Tx, gr GridRow) error
	SaveLineUpdate(tx *sql.Tx, gr GridRow) error
	SaveLineDelete(tx *sql.Tx, gr GridRow) error
}

type gridRowRepository struct {
	conn               *sql.DB
	tableName          string
	lineTableName      string
	parentFieldMapping map[string][]string
	childFieldMapping  map[string][]string
}

type SaveLineCallBack struct {
}

// NewGridRepository use for table pair with format table and table_lines, only for update
func NewGridRepository(conn *sql.DB, tableName, lineTableName string, parentFieldMapping, childFieldMapping map[string][]string) GridRowRepositoryWithChild {
	return &gridRowRepository{
		conn:               conn,
		tableName:          tableName,
		lineTableName:      lineTableName,
		parentFieldMapping: parentFieldMapping,
		childFieldMapping:  childFieldMapping,
	}
}

func (s *gridRowRepository) IsChild(gr GridRow) bool {
	id := gr.GetIDStr()

	return strings.HasSuffix(id, lineSuffix)
}

func (s *gridRowRepository) GetParentID(gr GridRow) (parentID interface{}, err error) {
	query := `
	SELECT parent_id 
	FROM ` + s.lineTableName + `
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, gr.GetID()).Scan(&parentID)

	return
}

// func (s *gridRowRepository) GetStatus(id interface{}) (status interface{}, err error) {
// 	query := `
// 	SELECT d.status
// 	FROM ` + s.tableName + ` t
// 		INNER JOIN documents d ON d.id = t.document_id
// 	WHERE t.id = ?
// 	`

// 	err = s.conn.QueryRow(query, id).Scan(&status)

// 	return
// }

func (s *gridRowRepository) Save(tx *sql.Tx, tr *MainRow) error {
	logger.Debug("Save grid row id", tr.IDString())

	if err := s.SaveMainRow(tx, tr); err != nil {
		return fmt.Errorf("save main row: [%w]", err)
	}

	if err := s.SaveLines(tx, tr); err != nil {
		return fmt.Errorf("save line: [%w]", err)
	}

	return nil
}

func (s *gridRowRepository) SaveMainRow(tx *sql.Tx, tr *MainRow) error {
	var (
		query string
		args  []interface{}
	)

	changedRow := GenGridRowChange(tr.Fields)

	switch tr.Fields.GetActionType() {
	case GridRowActionAdd:
		query, args = tr.Fields.MakeInsertQuery(s.tableName, s.parentFieldMapping)
		res, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("exec query: [%w], query: %s", err, query)
		}

		newID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("last inserted id: [%w]", err)
		}

		// if add success, set success color for row
		changedRow.Color = ChangedSuccessColor
		changedRow.Added = 1
		// set new id for row
		changedRow.NewId = fmt.Sprintf("%v$%d", changedRow.Parent, newID) // full id
		SetGridRowChangedResult(tr.Fields, changedRow)

		for k := range tr.Items {
			tr.Items[k]["Parent"] = newID

			// update parent id in ChangedRow of child row
			if childChangedRow, ok := tr.Items[k]["ChangedRow"].(ChangedRow); ok {
				childChangedRow.Parent = newID
				tr.Items[k]["ChangedRow"] = childChangedRow
			}
		}

		return nil
	case GridRowActionChanged:
		query, args = tr.Fields.MakeUpdateQuery(s.tableName, s.parentFieldMapping)
		args = append(args, tr.Fields.GetID())

		// parent contains only id - have nothing to update
		if len(args) == 1 {
			logger.Debug("Updates only id, nothing to update", tr.IDString(), args)

			return nil
		}

		if _, err := tx.Exec(query, args...); err != nil {
			return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
		}

		// if update success, set success color for row
		changedRow.Color = ChangedSuccessColor
		changedRow.Changed = 1
		SetGridRowChangedResult(tr.Fields, changedRow)

		return nil
	case GridRowActionDeleted:
		// ignore id start with CR
		idStr := tr.Fields.GetIDStr()
		if strings.HasPrefix(idStr, "CR") {
			logger.Debug("ignore this id: ", idStr)
			return nil
		}

		query, args = tr.Fields.MakeDeleteQuery(s.tableName)
		args = append(args, tr.Fields.GetID())

		if _, err := tx.Exec(query, args...); err != nil {
			return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
		}

		// if delete success, set success color for row
		changedRow.Color = ChangedSuccessColor
		changedRow.Deleted = 1
		SetGridRowChangedResult(tr.Fields, changedRow)

		return nil
	case GridRowActionNone:
		changedRow.Color = ""
		return nil
	default:
		return fmt.Errorf("undefined row type: %s", tr.Fields.GetActionType())
	}

}

func (s *gridRowRepository) SaveLines(tx *sql.Tx, tr *MainRow) error {
	for _, item := range tr.Items {
		var (
			query string
			args  []interface{}
		)

		changedRow := GenGridRowChange(item)

		switch item.GetActionType() {
		case GridRowActionAdd:
			// if parent is not persisted, no need to save child rows
			if !IsParentPersisted(item["Parent"]) {
				return fmt.Errorf("parent not saved")
			}

			query, args = item.MakeInsertQuery(s.lineTableName, s.childFieldMapping)

			res, err := tx.Exec(query, args...)
			if err != nil {
				return fmt.Errorf("exec query: [%w], query: %s", err, query)
			}

			newID, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("get last inserted id: [%w]", err)
			}

			// if add success, set success color for row
			changedRow.Color = ChangedSuccessColor
			changedRow.Added = 1
			// set new id for row
			changedRow.NewId = fmt.Sprintf("%v$%d%s", changedRow.Parent, newID, lineSuffix) // full id
			SetGridRowChangedResult(item, changedRow)

			continue
		case GridRowActionChanged:
			query, args = item.MakeUpdateQuery(s.lineTableName, s.childFieldMapping)
			args = append(args, getLineID(item.GetIDStr()))

			_, err := tx.Exec(query, args...)
			if err != nil {
				return fmt.Errorf("exec query: [%w], query: %s", err, query)
			}

			// if update success, set success color for row
			changedRow.Color = ChangedSuccessColor
			changedRow.Changed = 1
			SetGridRowChangedResult(item, changedRow)

			continue
		case GridRowActionDeleted:
			query, args = item.MakeDeleteQuery(s.lineTableName)
			args = append(args, getLineID(item.GetIDStr()))

			_, err := tx.Exec(query, args...)
			if err != nil {
				return fmt.Errorf("exec query: [%w], query: %s", err, query)
			}

			// if delete success, set success color for row
			changedRow.Color = ChangedSuccessColor
			changedRow.Deleted = 1
			SetGridRowChangedResult(item, changedRow)

			continue
		default:
			return fmt.Errorf("undefined row type: %s", item.GetActionType())
		}
	}

	return nil
}

func (s *gridRowRepository) SaveLineAdd(tx *sql.Tx, item GridRow) error {
	changedRow := GenGridRowChange(item)

	// if parent is not persisted, no need to save child rows
	if !IsParentPersisted(item["Parent"]) {
		return fmt.Errorf("parent not saved")
	}

	query, args := item.MakeInsertQuery(s.lineTableName, s.childFieldMapping)
	res, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("exec query: [%w], query: %s", err, query)
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last inserted id: [%w]", err)
	}

	// if add success, set success color for row
	changedRow.Color = ChangedSuccessColor
	changedRow.Added = 1
	// set new id for row
	changedRow.NewId = fmt.Sprintf("%v$%d%s", changedRow.Parent, newID, lineSuffix) // full id
	SetGridRowChangedResult(item, changedRow)
	return nil
}

func (s *gridRowRepository) SaveLineUpdate(tx *sql.Tx, item GridRow) error {
	changedRow := GenGridRowChange(item)

	query, args := item.MakeUpdateQuery(s.lineTableName, s.childFieldMapping)
	args = append(args, getLineID(item.GetIDStr()))

	_, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("exec query: [%w], query: %s", err, query)
	}

	// if update success, set success color for row
	changedRow.Color = ChangedSuccessColor
	changedRow.Changed = 1
	SetGridRowChangedResult(item, changedRow)

	return nil
}

func (s *gridRowRepository) SaveLineDelete(tx *sql.Tx, item GridRow) error {
	changedRow := GenGridRowChange(item)

	query, args := item.MakeDeleteQuery(s.lineTableName)
	args = append(args, getLineID(item.GetIDStr()))

	_, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("exec query: [%w], query: %s", err, query)
	}

	// if delete success, set success color for row
	changedRow.Color = ChangedSuccessColor
	changedRow.Deleted = 1
	SetGridRowChangedResult(item, changedRow)

	return nil
}

func getLineID(gridID string) (dbID string) {
	return strings.Trim(gridID, lineSuffix)
}
