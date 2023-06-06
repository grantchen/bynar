package treegrid

import (
	"fmt"
)

// because in case child row maybe has id which concat with parent's id: ex: 2-line has id: CR5$2-line, copy origin to mark this id
// this can be used when return result

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

// ParseRequestUploadSingleRow parse request upload from table without child
func ParseRequestUploadSingleRow(req *PostRequest) ([]GridRow, error) {
	grRowList := make([]GridRow, 0)

	for k := range req.Changes {
		ch := GridRow(req.Changes[k])
		ch.StoreGridTreeID()
		grRowList = append(grRowList, ch)
	}
	return grRowList, nil
}

// ParseRequestUpload parse request upload for parent-child table
func ParseRequestUpload(req *PostRequest, identityStore IdentityStorage) (*GridList, error) {
	trList := &GridList{
		mainRows:  make(map[string]GridRow),
		childRows: make(map[string][]GridRow),
	}

	// sort changes on transfers and items
	// logger.Debug("change len: ", len(req.Changes))
	for k := range req.Changes {
		ch := GridRow(req.Changes[k])
		// store origin id of gridtree, very useful in case row is child
		ch.StoreGridTreeID()

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
		trList.mainRows[k].StoreGridTreeID()
	}

	return trList, nil
}

// use by task2
func ParseRequestUpload2(req *PostRequest) (*GridList, error) {
	trList := &GridList{
		mainRows:  make(map[string]GridRow),
		childRows: make(map[string][]GridRow),
	}

	for k, change := range req.Changes {
		changeType, ok := change[ChangeTypeField].(string)
		if !ok {
			return nil, fmt.Errorf("undefined change type: [%v]", change[ChangeTypeField])
		}

		ID := change[FieldNameID].(string)

		// check if change is transfer of item
		switch changeType {
		case ChangeTypeNode:
			trList.mainRows[ID] = req.Changes[k]
		case ChangeTypeItem:
			parentID, _ := change[FieldNameParent].(string)

			if _, ok := trList.childRows[parentID]; !ok {
				trList.childRows[parentID] = make([]GridRow, 0, 10)
			}

			trList.childRows[parentID] = append(trList.childRows[parentID], req.Changes[k])

		default:
			return trList, fmt.Errorf("undefined change type: '%s'", changeType)
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
	mainRows := make([]*MainRow, 0, len(l.mainRows))

	for k := range l.mainRows {
		tr := &MainRow{
			Fields: l.mainRows[k],
			Items:  l.childRows[k],
		}

		mainRows = append(mainRows, tr)
	}

	return mainRows
}
