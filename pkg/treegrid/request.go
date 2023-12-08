package treegrid

import (
	"fmt"
	"log"
	"strconv"
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
		// set default ChangeRow for each row,
		// it will be set success if no error occurs in row processing(Add, Update, Delete)
		SetGridRowChangedResult(ch, GenGridRowChange(ch))
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
		// set default ChangeRow for each row,
		// it will be set success if no error occurs in row processing(Add, Update, Delete)
		SetGridRowChangedResult(ch, GenGridRowChange(ch))

		isChild, err := SetGridRowIdentity(ch, identityStore)
		if err != nil {
			return nil, fmt.Errorf("set gridRow identity: [%w]", err)
		}

		if !isChild {
			trList.mainRows[ch.GetIDStr()] = ch

			continue
		}

		// store parent id of child row
		ch.StoreGridParentID()
		parentID := ch.GetParentID()
		ch["Parent"] = parentID

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
			"id":                      k,
			string(GridRowActionNone): "1",
		}
		trList.mainRows[k].StoreGridTreeID()
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

// HandleSingleTreegridRows - handle rows of single treegrid one by one,
// and return response for Upload API
// save func is used to save each row
func HandleSingleTreegridRows(grList []GridRow, save func(gr GridRow) error) *PostResponse {
	resp := &PostResponse{Changes: []map[string]interface{}{}}
	for _, gr := range grList {
		if err := save(gr); err != nil {
			log.Println("Err", err)
			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
		}

		// set row change result(been set in row handle func)
		resp.Changes = append(resp.Changes, GetGridRowChangedResult(gr))
	}

	return resp
}

// HandleTreegridWithChildRows - handle rows of treegrid with child one by one,
// and return response for Upload API,
// saveMainRow is used to save parent row,
// saveLine is used to save child row
func HandleTreegridWithChildRows(trList *GridList, saveMainRow func(mr *MainRow) error, saveLine func(mr *MainRow, item GridRow) error) *PostResponse {
	resp := &PostResponse{}
	for _, mr := range trList.MainRows() {
		// handle parent row
		err := saveMainRow(mr)
		if err != nil {
			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
		}

		// if action is not none(is add, update or delete), set result color for parent row
		if mr.Fields.GetActionType() != GridRowActionNone {
			// set parent row change result(been set in row handle func)
			resp.Changes = append(resp.Changes, GetGridRowChangedResult(mr.Fields))
		}

		// failed to add new parent row, skip child rows handle
		// set error color(default RowChange) for child rows
		if err != nil && mr.Fields.GetActionType() == GridRowActionAdd {
			for _, item := range mr.Items {
				// default RowChange is error color
				resp.Changes = append(resp.Changes, GetGridRowChangedResult(item))
			}
			continue
		}

		for _, item := range mr.Items {
			// handle child row
			if err = saveLine(mr, item); err != nil {
				resp.IO.Result = -1
				resp.IO.Message += err.Error() + "\n"
			}

			// set child row change result(been set in row handle func)
			resp.Changes = append(resp.Changes, GetGridRowChangedResult(item))
		}
	}

	return resp
}

// HandleTreegridWithChildMainRowsLines - handle the main rows with lines of treegrid with child one by one,
// and return response for Upload API.
// All child rows under a parent row are in the same transaction as the parent row.
// saveMainRowWithLines is used to save parent row with child lines
func HandleTreegridWithChildMainRowsLines(trList *GridList, saveMainRowWithLines func(mr *MainRow) error) *PostResponse {
	resp := &PostResponse{}
	for _, mr := range trList.MainRows() {
		// handle parent row with child lines
		err := saveMainRowWithLines(mr)
		if err != nil {
			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"

			// if action is not none(is add, update or delete), set result color for parent row
			if mr.Fields.GetActionType() != GridRowActionNone {
				resp.Changes = append(resp.Changes, GenMapColorChangeError(mr.Fields))
			}

			for _, item := range mr.Items {
				// set child row error color
				resp.Changes = append(resp.Changes, GenMapColorChangeError(item))
			}
		} else {
			// set parent row change result(been set in row handle func)
			resp.Changes = append(resp.Changes, GetGridRowChangedResult(mr.Fields))

			for _, item := range mr.Items {
				// set child row change result(been set in row handle func)
				resp.Changes = append(resp.Changes, GetGridRowChangedResult(item))
			}
		}
	}

	return resp
}

// IsParentPersisted checks if parent is persisted(is a number id), default parent is letters id from front-end
func IsParentPersisted(parentId interface{}) bool {
	switch parentId.(type) {
	case int, int64:
		return true
	case string:
		parent, ok := parentId.(string)
		if !ok {
			return false
		}

		if _, err := strconv.Atoi(parent); err == nil {
			return true
		}
	}

	return false
}
