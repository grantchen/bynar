package treegrid

import (
	"fmt"
	"strings"
)

type Column struct {
	DBName      string
	DBNameShort string
	GridName    string
	IsItem      bool
	IsDate      bool
}

func NewColumn(gridName string, ChildFields map[string][]string, ParentFieldAliases map[string][]string) (c Column) {
	c.GridName = gridName
	if val, ok := ChildFields[gridName]; ok {
		c.IsItem = true
		c.DBName = val[0]
	}

	if val, ok := ParentFieldAliases[gridName]; ok {
		c.IsItem = false
		c.DBName = val[0]
	}

	if c.DBName == "" {
		c.DBName = gridName
	}

	nameParts := strings.Split(c.DBName, ".")
	if len(nameParts) > 0 {
		c.DBNameShort = nameParts[len(nameParts)-1]
	}

	if IsDateCol(c.DBName) {
		c.IsDate = true
	}

	return
}

// WhereSQL return where sql for column
func (c Column) WhereSQL(val string) (where string, arg interface{}) {
	if c.IsDate {
		if val == "" {
			return fmt.Sprintf("%s IS NULL", c.DBName), nil
		}
	}

	return fmt.Sprintf("%s = ?", c.DBName), val
}

// Type return type for column
func (c Column) Type() string {
	if c.IsDate {
		return "Date"
	}

	if IsIntCol(c.DBName) {
		return "Int"
	}

	return "Text"
}
