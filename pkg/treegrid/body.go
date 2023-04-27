package treegrid

import "strconv"

type BodyParam struct {
	ID   string `json:"id,omitempty"`
	Rows string `json:"rows,omitempty"`
	Pos  string `json:"pos,omitempty"`
	Col  string `json:"Col,omitempty"`
	Val  string `json:"Val,omitempty"`
}

// / GetRowLevel gets from rows level. Example: rows = "2WHERE filed=1", level = 1
func (b BodyParam) GetRowLevel() int {
	if b.Rows == "" {
		return 0
	}

	id, _ := strconv.Atoi(b.Rows[:1])

	return id
}

// GetRowWhere gets from rows where clause. Example: rows = "2WHERE filed=1", return WHERE filed=1
func (b BodyParam) GetRowWhere() string {
	if b.Rows == "" {
		return ""
	}

	return b.Rows[1:]
}

// GetItemsRequest defines type of response - Transfer of TransferItems
// Conditions for items response:  "id" is digit and "rows" == ""
func (b BodyParam) GetItemsRequest() bool {
	if _, ok := b.IntID(); !ok {
		return false
	}

	if b.Rows != "" {
		return false
	}

	return true
}

// IntID converts id string to int with indicating id existence
func (b BodyParam) IntID() (int, bool) {
	if b.ID == "" {
		return 0, false
	}

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
