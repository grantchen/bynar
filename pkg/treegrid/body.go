package treegrid

import (
	"fmt"
	"strconv"
	"strings"
)

type BodyParam struct {
	ID               string `json:"id,omitempty"`
	Rows             string `json:"rows,omitempty"`
	Pos              string `json:"pos,omitempty"`
	Col              string `json:"Col,omitempty"`
	Val              string `json:"Val,omitempty"`
	TreegridOriginID string // using for keep id from treegrid, especially with child row, and group with $
}

// SetBodyParamRows sets rows param for body. Example: rows = "2|13|19|WHERE filed=1WHERE child_field=2"
func SetBodyParamRows(level int, parentWhere, childWhere string) string {
	return fmt.Sprintf("%d|%d|%d|%s%s", level, len(parentWhere), len(childWhere), parentWhere, childWhere)
}

// GetRowLevel gets from rows level. Example: rows = "2WHERE filed=1", level = 1
func (b *BodyParam) GetRowLevel() int {
	if b.Rows == "" {
		return 0
	}

	id, _ := strconv.Atoi(b.Rows[:1])

	return id
}

// GetRowParentWhere gets from rows parent where clause. Example: rows = "2|13|19|WHERE filed=1WHERE child_field=2"
func (b *BodyParam) GetRowParentWhere() string {
	if b.Rows == "" {
		return ""
	}

	splitN := strings.SplitN(b.Rows, "|", 4)
	if len(splitN) != 4 {
		return ""
	}

	parentLen, _ := strconv.Atoi(splitN[1])
	return splitN[3][:parentLen]
}

// GetRowChildWhere gets from rows child where clause. Example: rows = "2|13|19|WHERE filed=1WHERE child_field=2"
func (b *BodyParam) GetRowChildWhere() string {
	if b.Rows == "" {
		return ""
	}

	splitN := strings.SplitN(b.Rows, "|", 4)
	if len(splitN) != 4 {
		return ""
	}

	parentLen, _ := strconv.Atoi(splitN[1])
	return splitN[3][parentLen:]
}

// GetItemsRequest defines type of response - Transfer of TransferItems
// Conditions for items response:  "id" is digit and "rows" == ""
func (b *BodyParam) GetItemsRequest() bool {

	//check is group, move outside, parse request
	// if strings.Contains(b.ID, "$") {
	// 	idGroup := strings.Split(b.ID, "$")
	// 	newId := idGroup[len(idGroup)-1]
	// 	b.ID = newId
	// }

	if _, ok := b.IntID(); !ok {
		return false
	}

	if b.Rows != "" {
		return false
	}

	return true
}

// IntID converts id string to int with indicating id existence
func (b *BodyParam) IntID() (int, bool) {
	if b.ID == "" {
		return 0, false
	}

	//check is group, move to outside, parse request
	// if strings.Contains(b.ID, "$") {
	// 	idGroup := strings.Split(b.ID, "$")
	// 	newId := idGroup[len(idGroup)-1]
	// 	b.ID = newId
	// }

	id, err := strconv.Atoi(b.ID)
	if err != nil {
		return 0, false
	}

	return id, true
}

// IntID converts id string to int with indicating id existence
func (b BodyParam) IntPos() (int, bool) {
	if b.Pos == "" {
		return 0, false
	}

	id, err := strconv.Atoi(b.Pos)
	if err != nil {
		return 0, false
	}

	return id, true
}
