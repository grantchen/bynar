package treegrid

import "strings"

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
		c.DBName = val[0]
	}

	if c.DBName == "" {
		c.DBName = gridName
	}

	nameParts := strings.Split(c.DBName, ".")
	if len(nameParts) > 0 {
		c.DBNameShort = nameParts[len(nameParts)-1]
	}

	return
}
