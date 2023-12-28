package treegrid

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type BodyParam struct {
	ID               string `json:"id,omitempty"`
	Rows             string `json:"rows,omitempty"`
	Pos              string `json:"pos,omitempty"`
	Col              string `json:"Col,omitempty"`
	Val              string `json:"Val,omitempty"`
	TreegridOriginID string // using for keep id from treegrid, especially with child row, and group with $

	RowsFields []RowsFieldCond // using for keep rows
	RowsLevel  int             // using for keep rows level
}

// RowsFieldCond defines rows field
type RowsFieldCond struct {
	GridName string `json:"name"`
	Value    string `json:"value"`
}

// ParseRows parses rows string param
func (b *BodyParam) ParseRows() error {
	if b.Rows == "" {
		b.RowsFields = make([]RowsFieldCond, 0)
		return nil
	}

	var err error
	// parse rows level
	b.RowsLevel, err = strconv.Atoi(b.Rows[:1])
	if err != nil {
		return err
	}

	// parse rows map
	var rowsM []RowsFieldCond
	if err = json.Unmarshal([]byte(b.Rows[1:]), &rowsM); err != nil {
		return err
	}
	b.RowsFields = rowsM

	return nil
}

// GetNewRows gets new rows string param
func (b *BodyParam) GetNewRows(level int, cond RowsFieldCond) (string, error) {
	newRowsFields := b.RowsFields
	newRowsFields = append(newRowsFields, cond)
	paramsJSON, err := json.Marshal(newRowsFields)
	if err != nil {
		return "", err
	}
	rows := fmt.Sprintf("%d%s", level, string(paramsJSON))
	return rows, nil
}

// GetItemsRequest defines type of response - Transfer of TransferItems
// Conditions for items response:  "id" is digit and "rows" == ""
func (b *BodyParam) GetItemsRequest() bool {
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

	id, err := strconv.Atoi(b.ID)
	if err != nil {
		return 0, false
	}

	return id, true
}

// IntPos converts id string to int with indicating id existence
func (b *BodyParam) IntPos() (int, bool) {
	if b.Pos == "" {
		return 0, false
	}

	id, err := strconv.Atoi(b.Pos)
	if err != nil {
		return 0, false
	}

	return id, true
}
