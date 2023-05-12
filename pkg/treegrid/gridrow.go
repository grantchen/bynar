package treegrid

import (
	"fmt"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
)

type (
	GridRowActionType string
	GridRow           map[string]interface{}
)

const (
	GridRowActionAdd     GridRowActionType = "Added"
	GridRowActionDeleted GridRowActionType = "Deleted"
	GridRowActionChanged GridRowActionType = "Changed"
)

func (f GridRow) IsChild() bool {
	panic("not implemented yet")
}

func (f GridRow) GetParentID() string {
	pID, _ := f.GetValString("Parent")

	return pID
}

func (f GridRow) ValidateOnRequiredAll(fieldsMapping map[string][]string) error {
	for key, _ := range fieldsMapping {
		_, ok := f[key]
		if !ok {
			return fmt.Errorf("[%w]: %s", errors.ErrMissingRequiredParams, "key")
		}
	}
	return nil
}

// MakeInsertQuery - returns query and args for query execution
// fieldsMapping - mapping for converting treegrid name to db name
// example: intput - map["user_group_id"] = []string {"group_id", "user.group_id"}.
func (f GridRow) MakeInsertQuery(tableName string, fieldsMapping map[string][]string) (query string, args []interface{}) {
	columnNames := make([]string, 0, len(f))
	args = make([]interface{}, 0, len(f))

	for treegridName, val := range f {
		if treegridName == "id" {
			continue
		}

		if treegridName == "transaction_no" {
			continue
		}

		if strVal, ok := val.(string); ok && strVal == "" {
			continue
		}

		if dbNames, ok := fieldsMapping[treegridName]; ok {
			columnNames = append(columnNames, dbNames[0])
			args = append(args, f[treegridName])
		}
	}

	if len(fieldsMapping["transaction_no"]) > 0 {
		columnNames = append(columnNames, "transaction_no")
	}

	vals := ""
	colNamesStr := strings.Join(columnNames, ",")
	for _, v := range columnNames {
		if v == "transaction_no" {
			vals += "UUID_SHORT(),"
			continue
		}

		if !strings.Contains(v, "_date") {
			vals += "?,"
			continue
		}
		vals += `STR_TO_DATE( ?, "%m/%d/%Y"),`
	}

	if len(vals) > 0 {
		vals = vals[:len(vals)-1]
	}

	query = fmt.Sprintf(`
	INSERT INTO %s (%s)
	VALUES (%s)
	`, tableName, colNamesStr, vals)

	return
}

// MakeUpdateQuery - returns query and args for query execution
// fieldsMapping - mapping for converting treegrid name to db name
// example: intput - map["user_group_id"] = []string {"group_id", "user.group_id"}.
func (f GridRow) MakeUpdateQuery(tableName string, fieldsMapping map[string][]string) (query string, args []interface{}) {
	query = `
	UPDATE %s
	SET %s
	WHERE id = ?
	`

	var (
		setsQuery string
	)
	args = make([]interface{}, 0, len(f))

	for treegridName := range f {
		if treegridName == "id" {
			continue
		}

		if dbNames, ok := fieldsMapping[treegridName]; ok {
			setsQuery += dbNames[0] + " = ?,"
			args = append(args, f[treegridName])
		}
	}

	return fmt.Sprintf(query, tableName, strings.Trim(setsQuery, ",")), args
}

func (f GridRow) MakeDeleteQuery(tableName string) (query string, args []interface{}) {
	args = make([]interface{}, 0)

	query = `
	DELETE FROM %s
	WHERE id = ?
	`

	return fmt.Sprintf(query, tableName), args
}

func (g GridRow) GetActionType() GridRowActionType {
	for name := range g {
		switch name {
		case string(GridRowActionAdd):
			return GridRowActionAdd
		case string(GridRowActionDeleted):
			return GridRowActionDeleted
		case string(GridRowActionChanged):
			return GridRowActionChanged
		}
	}

	return GridRowActionAdd
}

func (g GridRow) GetIDStr() (id string) {
	if val, ok := g["NewId"]; ok {
		return val.(string)
	}

	for name, val := range g {
		if name == "id" {
			id, _ = val.(string)
		}
	}

	return
}

func (g GridRow) GetIDInt() (id int) {
	id, _ = strconv.Atoi(g.GetIDStr())
	return id
}

func (g GridRow) GetID() (id interface{}) {
	if val, ok := g["NewId"]; ok {
		return val
	}

	for name, val := range g {
		if name == "id" {
			id = val

			return
		}
	}

	return
}

func (g GridRow) GetValString(name string) (string, bool) {
	val, ok := g[name]
	if !ok {
		return "", false
	}

	switch v := val.(type) {
	case int:
		return strconv.Itoa(v), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case string:
		return v, true
	}

	return "", false
}

func (g GridRow) GetValInt(name string) (int, bool) {
	val, ok := g[name]
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case string:
		intVal, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}

		return intVal, true
	}

	return 0, false
}

func (g GridRow) UpdatedFields() []string {
	updatedFields := make([]string, 0, len(g))
	for key := range g {
		switch key {
		case string(GridRowActionAdd), string(GridRowActionChanged), string(GridRowActionDeleted):
			continue
		case "Def":
			continue
		case "id":
			continue
		default:
			updatedFields = append(updatedFields, key)
		}
	}

	return updatedFields
}

func (g GridRow) GetStrInt(name string) (int, bool) {
	val, ok := g[name]
	if !ok {
		return 0, false
	}

	valInt, ok := val.(int)
	if !ok {
		return 0, false
	}

	return valInt, true
}
