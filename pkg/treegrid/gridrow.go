package treegrid

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"strconv"
	"strings"
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

// use when grouping by
const reqID = "reqID"
const reqParentID = "reqParentID"

func (f GridRow) IsChild() bool {
	panic("not implemented yet")
}

func (f GridRow) GetParentID() string {
	return f.getRealID(f.getOriginParentID())
}

func (f GridRow) GetLineID() string {
	return strings.Trim(f.GetIDStr(), lineSuffix)
}

func (f GridRow) ValidateOnRequiredAll(fieldsMapping map[string][]string, language string) error {
	for key, _ := range fieldsMapping {
		val, ok := f[key]
		if !ok || val == "" {
			templateData := map[string]string{
				"Field": key,
			}
			return i18n.TranslationI18n(language, "RequiredFieldsBlank", templateData)
		}
	}
	return nil
}

// used to check empty update key.
func (f GridRow) ValidateOnRequired(fieldsMapping map[string][]string, language string) error {
	for key, _ := range fieldsMapping {
		if key == "Changed" || key == "id" {
			continue
		}
		val, ok := f[key]
		if ok && val == "" {
			templateData := map[string]string{
				"Field": key,
			}
			return i18n.TranslationI18n(language, "RequiredFieldsBlank", templateData)
		}
	}
	return nil
}

// Check Field Limit length.
func (f GridRow) ValidateOnLimitLength(fieldsMapping map[string][]string, limitLength int, language string) error {
	for key, _ := range fieldsMapping {
		if key == "Changed" || key == "id" {
			continue
		}
		_, ok := f[key]
		stringVal, isString := f.GetValString(key)
		if ok && isString && len(stringVal) > limitLength {
			templateData := map[string]string{
				"Field": key,
			}
			return i18n.TranslationI18n(language, "FieldTooLong", templateData)
		}
	}
	return nil
}

// Check Field Limit length to float.
func (f GridRow) ValidateOnLimitLengthToFloat(fieldsMapping map[string][]string, language string) error {
	for key, _ := range fieldsMapping {
		if key == "Changed" || key == "id" {
			continue
		}
		_, ok := f[key]
		templateData := map[string]string{
			"Field": key,
		}
		numberVal, _ := f.GetValFloat64(key)
		// Database default write limit 10000000000 to Float
		if ok && numberVal > 10000000000 {
			return i18n.TranslationI18n(language, "FieldOutRange", templateData)
		}
		stringVal, isString := f.GetValString(key)
		if ok && isString {
			// treegrid input number automatically add spaces after inputting more than 10 characters 'e '
			if strings.Contains(stringVal, "e ") {
				return i18n.TranslationI18n(language, "FieldOutRange", templateData)
			}
		}
	}
	return nil
}

// used to check not negative number.
func (f GridRow) ValidateOnNotNegativeNumber(fieldsMapping map[string][]string, language string) error {
	for key, _ := range fieldsMapping {
		if key == "Changed" || key == "id" {
			continue
		}
		_, ok := f[key]
		numberVal, _ := f.GetValFloat64(key)
		if ok && numberVal < 0 {
			templateData := map[string]string{
				"Field": key,
			}
			return i18n.TranslationI18n(language, "ValidateOnNotNegativeNumber", templateData)
		}

	}
	return nil
}

// used to check A positive number.
func (f GridRow) ValidateOnPositiveNumber(fieldsMapping map[string][]string, language string) error {
	for key, _ := range fieldsMapping {
		if key == "Changed" || key == "id" {
			continue
		}
		_, ok := f[key]
		numberVal, _ := f.GetValFloat64(key)
		if ok && numberVal <= 0 {
			templateData := map[string]string{
				"Field": key,
			}
			return i18n.TranslationI18n(language, "ValidateOnPositiveNumber", templateData)
		}

	}
	return nil
}

func (f GridRow) MakeValidateOnIntegrityQuery(tableName string, fieldsMapping map[string][]string, fieldsValidating []string) (query string, args []interface{}) {
	queryFormat := `select COUNT(*) as Count
	FROM %s
	WHERE 1=1 %s `

	var whereCondition string
	args = make([]interface{}, 0)
	for _, field := range fieldsValidating {
		dbFields, ok := fieldsMapping[field]
		if !ok {
			continue
		}
		whereCondition += fmt.Sprintf(" AND %s = ? ", dbFields[0])
		args = append(args, f[field])
	}

	// only when action is changed, we need to add condition id != ?
	if f.GetActionType() == GridRowActionChanged {
		whereCondition += " AND id != ? "
		args = append(args, f.GetID())
	}

	query = fmt.Sprintf(queryFormat, tableName, whereCondition)
	return
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

		//if treegridName == "transaction_no" {
		//	continue
		//}

		if dbNames, ok := fieldsMapping[treegridName]; ok {
			columnName := dbNames[0]
			if strVal, ok := val.(string); ok && strVal == "" {
				// date column should insert with null
				if IsDateCol(columnName) {
					continue
				}
			}

			columnNames = append(columnNames, columnName)
			args = append(args, f[treegridName])
		}
	}

	//if len(fieldsMapping["transaction_no"]) > 0 {
	//	columnNames = append(columnNames, "transaction_no")
	//}

	vals := ""
	colNamesStr := strings.Join(columnNames, ",")
	for _, v := range columnNames {
		//if v == "transaction_no" {
		//	vals += "UUID_SHORT(),"
		//	continue
		//}

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
			columnName := dbNames[0]
			if strings.Contains(columnName, "_date") {
				// If the field is a time type, convert the time string to a time type
				setsQuery += columnName + " = STR_TO_DATE(?, '%m/%d/%Y'),"
			} else {
				setsQuery += columnName + " = ?,"
			}
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
		return g.getRealID(val.(string))
	}

	for name, val := range g {
		if name == "id" {
			id, _ = val.(string)
			id = g.getRealID(id)
		}
	}

	return
}

func (g GridRow) GetIDInt() (id int) {
	id, _ = strconv.Atoi(g.GetIDStr())
	return id
}

// getRealID return real id of row
func (g *GridRow) getRealID(id string) string {
	if strings.Contains(id, "$") { // splitId when group by have format: (CR[0-9]+\$)+<real_id>
		realId := ""
		for _, splitId := range strings.Split(id, "$") {
			if !g.isAutoID(splitId) {
				// real id is the last
				realId = splitId
			}
		}
		return realId
	}

	return id
}

// isAutoID check if id is auto generated id
// AutoIdPrefix="AR" GroupIdPrefix="GR" ChildIdPrefix="CR"
func (g *GridRow) isAutoID(id string) bool {
	if len(id) < 2 {
		return false
	}

	idPrefix := id[:2]
	return idPrefix == "AR" || idPrefix == "GR" || idPrefix == "CR"
}

// getRealIDInterface return real id interface of row
func (g *GridRow) getRealIDInterface(input interface{}) interface{} {
	idStr, _ := input.(string)
	if strings.Contains(idStr, "$") { // splitId when group by have format: (CR[0-9]+\$)+<real_id>
		realId := ""
		for _, splitId := range strings.Split(idStr, "$") {
			if !g.isAutoID(splitId) {
				// real id is the last
				realId = splitId
			}
		}
		return realId
	}
	return input
}

func (g GridRow) GetID() (id interface{}) {
	return g.getRealIDInterface(g.getOriginID())
}

func (g GridRow) getOriginID() (id interface{}) {
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

// return raw id from treegrid req, used for grouping feature when id include parent: ex id: 2-line => CR5$2-line
func (g GridRow) GetTreeGridID() (id interface{}) {
	id, ok := g[reqID]
	if ok {
		return id
	}
	// logger.Debug("orgin id: ", g[reqID])
	return g.GetID()
}

func (g GridRow) StoreGridTreeID() {
	g[reqID] = g.getOriginID()
}

// StoreGridParentID store parent id of child row
func (g GridRow) StoreGridParentID() {
	g[reqParentID] = g.getOriginParentID()
}

// return Parent from treegrid req
func (g GridRow) getOriginParentID() string {
	pID, _ := g.GetValString("Parent")
	return pID
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

// GetValFloat64 return float64 value of field name
func (g GridRow) GetValFloat64(name string) (float64, bool) {
	val, ok := g[name]
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case string:
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return floatVal, true
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
		case "reqID":
			continue
		case "Custom":
			continue
		default:
			updatedFields = append(updatedFields, key)
		}
	}

	return updatedFields
}

func (g GridRow) MergeWithMap(m map[string]interface{}) GridRow {
	newG := make(map[string]interface{})

	for k, v := range g {
		newG[k] = v
	}

	for k, v := range m {
		newG[k] = v
	}

	return newG
}

// FilterFieldsMapping return new fieldsMapping with only fields in fields
func (f GridRow) FilterFieldsMapping(fieldsMapping map[string][]string, fields []string) map[string][]string {
	newFieldsMapping := make(map[string][]string)
	for _, field := range fields {
		if _, ok := fieldsMapping[field]; ok {
			newFieldsMapping[field] = fieldsMapping[field]
		}
	}
	return newFieldsMapping
}
