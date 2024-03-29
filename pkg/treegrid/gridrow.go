package treegrid

import (
	"fmt"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
)

type (
	GridRowActionType string
	GridRow           map[string]interface{}
)

const (
	GridRowActionAdd     GridRowActionType = "Added"
	GridRowActionDeleted GridRowActionType = "Deleted"
	GridRowActionChanged GridRowActionType = "Changed"

	// GridRowActionNone is no need to handle tag,
	// if only child lines change without a parent, no need to handle parent.
	// ParseRequestUpload in pkg/treegrid/request.go will create a parent row as parent of child lines
	// and set the action type of row to GridRowActionNone, means no need to handle parent row
	GridRowActionNone GridRowActionType = "None" // no need to handle
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
	for key := range fieldsMapping {
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

// ValidateOnRequired used to check empty update key.
func (f GridRow) ValidateOnRequired(fieldsMapping map[string][]string, language string) error {
	for key := range fieldsMapping {
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

// ValidateOnLimitLength Check Field Limit length.
func (f GridRow) ValidateOnLimitLength(fieldsMapping map[string][]string, limitLength int, language string) error {
	for key := range fieldsMapping {
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

// ValidateOnLimitLengthToFloat Check Field Limit length to float.
func (f GridRow) ValidateOnLimitLengthToFloat(fieldsMapping map[string][]string, language string) error {
	for key := range fieldsMapping {
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

// ValidateOnNotNegativeNumber used to check not negative number.
func (f GridRow) ValidateOnNotNegativeNumber(fieldsMapping map[string][]string, language string) error {
	for key := range fieldsMapping {
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

// ValidateOnPositiveNumber used to check A positive number.
func (f GridRow) ValidateOnPositiveNumber(fieldsMapping map[string][]string, language string) error {
	for key := range fieldsMapping {
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

func (f GridRow) GetActionType() GridRowActionType {
	// if a new added row is deleted, it should be deleted
	// it has "Added" and "Deleted" tag, example: {"id":"_","Deleted":1,"Added":1}]}
	if _, ok := f["Deleted"]; ok {
		return GridRowActionDeleted
	}

	for name := range f {
		switch name {
		case string(GridRowActionAdd):
			return GridRowActionAdd
		case string(GridRowActionChanged):
			return GridRowActionChanged
		case string(GridRowActionNone):
			return GridRowActionNone
		}
	}

	return GridRowActionAdd
}

func (f GridRow) GetIDStr() (id string) {
	if val, ok := f["NewId"]; ok {
		return f.getRealID(val.(string))
	}

	for name, val := range f {
		if name == "id" {
			id, _ = val.(string)
			id = f.getRealID(id)
		}
	}

	return
}

func (f GridRow) GetIDInt() (id int) {
	id, _ = strconv.Atoi(f.GetIDStr())
	return id
}

// getRealID return real id of row
func (f GridRow) getRealID(id string) string {
	if strings.Contains(id, "$") { // splitId when group by have format: (CR[0-9]+\$)+<real_id>
		realId := ""
		for _, splitId := range strings.Split(id, "$") {
			if !f.isAutoID(splitId) {
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
func (f GridRow) isAutoID(id string) bool {
	if len(id) < 2 {
		return false
	}

	idPrefix := id[:2]
	return idPrefix == "AR" || idPrefix == "GR" || idPrefix == "CR"
}

// getRealIDInterface return real id interface of row
func (f GridRow) getRealIDInterface(input interface{}) interface{} {
	idStr, _ := input.(string)
	if strings.Contains(idStr, "$") { // splitId when group by have format: (CR[0-9]+\$)+<real_id>
		realId := ""
		for _, splitId := range strings.Split(idStr, "$") {
			if !f.isAutoID(splitId) {
				// real id is the last
				realId = splitId
			}
		}
		return realId
	}
	return input
}

func (f GridRow) GetID() (id interface{}) {
	return f.getRealIDInterface(f.getOriginID())
}

func (f GridRow) getOriginID() (id interface{}) {
	if val, ok := f["NewId"]; ok {
		return val
	}

	for name, val := range f {
		if name == "id" {
			id = val

			return
		}
	}

	return
}

// GetTreeGridID return raw id from treegrid req, used for grouping feature when id include parent: ex id: 2-line => CR5$2-line
func (f GridRow) GetTreeGridID() (id interface{}) {
	id, ok := f[reqID]
	if ok {
		return id
	}
	// logger.Debug("orgin id: ", g[reqID])
	return f.GetID()
}

func (f GridRow) StoreGridTreeID() {
	f[reqID] = f.getOriginID()
}

// StoreGridParentID store parent id of child row
func (f GridRow) StoreGridParentID() {
	f[reqParentID] = f.getOriginParentID()
}

// return Parent from treegrid req
func (f GridRow) getOriginParentID() string {
	pID, _ := f.GetValString("Parent")
	return pID
}

func (f GridRow) GetValString(name string) (string, bool) {
	val, ok := f[name]
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

func (f GridRow) GetValInt(name string) (int, bool) {
	val, ok := f[name]
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
func (f GridRow) GetValFloat64(name string) (float64, bool) {
	val, ok := f[name]
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

func (f GridRow) UpdatedFields() []string {
	updatedFields := make([]string, 0, len(f))
	for key := range f {
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

func (f GridRow) MergeWithMap(m map[string]interface{}) GridRow {
	newG := make(map[string]interface{})

	for k, v := range f {
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
