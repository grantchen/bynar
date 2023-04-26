package treegrid

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type (
	GridList struct {
		mainRows  map[string]GridRow
		childRows map[string][]GridRow
	}

	IdentityStorage interface {
		GetParentID(gr GridRow) (parentID interface{}, err error)
		IsChild(gr GridRow) bool
	}
)

func ParseRequest(req *treegrid.PostRequest, identityStore IdentityStorage) (*GridList, error) {
	trList := &GridList{
		mainRows:  make(map[string]GridRow),
		childRows: make(map[string][]GridRow),
	}

	// sort changes on transfers and items
	for k := range req.Changes {
		ch := GridRow(req.Changes[k])

		isChild, err := SetGridRowIdentity(ch, identityStore)
		if err != nil {
			return nil, fmt.Errorf("set gridRow identity: [%w]", err)
		}

		if !isChild {
			trList.mainRows[ch.GetIDStr()] = ch

			continue
		}

		parentID := ch.GetParentID()
		if _, ok := trList.childRows[parentID]; !ok {
			trList.childRows[parentID] = make([]GridRow, 0, 10)
		}

		trList.childRows[parentID] = append(trList.childRows[parentID], ch)
	}

	// if get only items without transfer object
	// create transfer object with only id field
	// k = parentID
	// rows - items for transfer with ID = k
	for k := range trList.childRows {
		_, ok := trList.mainRows[k]
		if ok {
			continue
		}

		trList.mainRows[k] = GridRow{
			"id":                         k,
			string(GridRowActionChanged): "1",
		}
	}

	return trList, nil
}

// SetGridRowIdentity sets required params for indentifying grid row
// sets params:
// "Def": "Node" | "Data"
func SetGridRowIdentity(gr GridRow, identityStore IdentityStorage) (isChild bool, err error) {
	id := gr.GetID()

	// sets attributes based on "Def" field
	if val, ok := gr.GetValString("Def"); ok {
		// Nothing need to set
		if val == "Node" {
			return false, nil
		}

		if val == "Data" { // all data already are set
			return true, nil
		}

		return false, fmt.Errorf("undefined 'Def' value: %s", val)
	}

	if identityStore.IsChild(gr) {
		parentID, err := identityStore.GetParentID(gr)
		if err != nil {
			return true, fmt.Errorf("get parent id: [%w], by line id: %v", err, id)
		}

		gr["Parent"] = parentID

		return true, nil
	}

	return
}

// MainRows - contains all received transfer attributes with related items
func (l *GridList) MainRows() []*MainRow {
	transfers := make([]*MainRow, 0, len(l.mainRows))

	for k := range l.mainRows {
		tr := &MainRow{
			Fields: l.mainRows[k],
			Items:  l.childRows[k],
		}

		transfers = append(transfers, tr)
	}

	return transfers
}
