package treegrid

import (
	"fmt"
	"strconv"
)

type (
	MainRow struct {
		Fields GridRow
		Items  []GridRow
	}
)

func (t *MainRow) IDString() string {
	idStr, ok := t.Fields["id"].(string)
	if !ok {
		return fmt.Sprintf("%v", t.Fields["id"])
	}

	return idStr
}

func (t *MainRow) Status() int {
	switch v := t.Fields["status"].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case string:
		id, _ := strconv.Atoi(v)

		return id
	default:
		return 0
	}
}
