package gridrep

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

const (
	lineSuffix = "-line"
)

type gridRowReppository struct {
	conn               *sql.DB
	tableName          string
	lineTableName      string
	parentFieldMapping map[string][]string
	childFieldMapping  map[string][]string
}

func NewGridRepository(conn *sql.DB, tableName, lineTableName string, parentFieldMapping, childFieldMapping map[string][]string) GridRowReppository {
	return &gridRowReppository{
		conn:               conn,
		tableName:          tableName,
		lineTableName:      lineTableName,
		parentFieldMapping: parentFieldMapping,
		childFieldMapping:  childFieldMapping,
	}
}

func (s *gridRowReppository) IsChild(gr treegrid.GridRow) bool {
	id := gr.GetIDStr()

	return strings.HasSuffix(id, lineSuffix)
}

func (s *gridRowReppository) GetParentID(gr treegrid.GridRow) (parentID interface{}, err error) {
	query := `
	SELECT parent_id 
	FROM ` + s.lineTableName + `
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, gr.GetID()).Scan(&parentID)

	return
}

func (s *gridRowReppository) GetStatus(id interface{}) (status interface{}, err error) {
	query := `
	SELECT d.status 
	FROM ` + s.tableName + ` t
		INNER JOIN documents d ON d.id = t.document_id
	WHERE t.id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&status)

	return
}

func (s *gridRowReppository) Save(tx *sql.Tx, tr *treegrid.MainRow) error {
	logger.Debug("Save grid row id", tr.IDString())

	if err := s.saveMainRow(tx, tr); err != nil {
		return fmt.Errorf("save main row: [%w]", err)
	}

	if err := s.saveLines(tx, tr); err != nil {
		return fmt.Errorf("save line: [%w]", err)
	}

	return nil
}

func (s *gridRowReppository) saveMainRow(tx *sql.Tx, tr *treegrid.MainRow) error {
	var (
		query string
		args  []interface{}
	)

	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		logger.Debug("add new row")

		query, args = tr.Fields.MakeInsertQuery(s.tableName, s.parentFieldMapping)

		logger.Debug(query, "args", args)
		res, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("exec query: [%w], query: %s", err, query)
		}

		newID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("last inserted id: [%w]", err)
		}

		// update id for row and child items
		tr.Fields["NewId"] = newID
		tr.Fields["Added"] = "1"

		logger.Debug("New row id", newID)
		logger.Debug("Update items field 'Parent'")

		for k := range tr.Items {
			tr.Items[k]["Parent"] = newID
		}

		return nil
	case treegrid.GridRowActionChanged:
		logger.Debug("update parent row")

		query, args = tr.Fields.MakeUpdateQuery(s.tableName, s.parentFieldMapping)
		args = append(args, tr.Fields.GetID())

		// parent contains only id - have nothing to update
		if len(args) == 1 {
			logger.Debug("Updates only id, nothing to update", tr.IDString())

			return nil
		}

		if _, err := tx.Exec(query, args...); err != nil {
			return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
		}

		return nil
	case treegrid.GridRowActionDeleted:
		logger.Debug("delete parent")

		query, args = tr.Fields.MakeDeleteQuery(s.tableName)
		args = append(args, tr.Fields.GetID())
	default:
		return fmt.Errorf("undefined row type: %s", tr.Fields.GetActionType())
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
	}

	return nil
}

func (s *gridRowReppository) saveLines(tx *sql.Tx, tr *treegrid.MainRow) error {
	logger.Debug("Save lines, count: ", len(tr.Items))

	for _, item := range tr.Items {
		var (
			query string
			args  []interface{}
		)

		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			query, args = item.MakeInsertQuery(s.lineTableName, s.childFieldMapping)

			res, err := tx.Exec(query, args...)
			if err != nil {
				return fmt.Errorf("exec query: [%w], query: %s", err, query)
			}

			newID, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("get last inserted id: [%w]", err)
			}

			item["NewId"] = fmt.Sprintf("%d%s", newID, lineSuffix)

			continue
		case treegrid.GridRowActionChanged:
			query, args = item.MakeUpdateQuery(s.lineTableName, s.childFieldMapping)
			args = append(args, getLineID(item.GetIDStr()))

			_, err := tx.Exec(query, args...)
			if err != nil {
				return fmt.Errorf("exec query: [%w], query: %s", err, query)
			}

			continue
		case treegrid.GridRowActionDeleted:
			query, args = item.MakeDeleteQuery(s.lineTableName)
			args = append(args, getLineID(item.GetIDStr()))
		default:
			return fmt.Errorf("undefined row type: %s", item.GetActionType())
		}

		_, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("exec query: [%w], query: %s", err, query)
		}
	}

	return nil
}

func getLineID(gridID string) (dbID string) {
	return strings.Trim(gridID, lineSuffix)
}
