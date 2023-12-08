package treegrid

import "encoding/json"

type (
	// ChangedRow: used to return Messages for POST update
	ChangedRow struct {
		Id      interface{} `json:"id,omitempty"`
		NewId   string      `json:"NewId,omitempty"`
		Changed int         `json:"Changed,omitempty"`
		Added   int         `json:"Added,omitempty"`
		Deleted int         `json:"Deleted,omitempty"`
		Color   string      `json:"Color,omitempty"`

		// for keeping row position
		Parent interface{} `json:"Parent,omitempty"` // Parent row id
		Next   interface{} `json:"Next,omitempty"`   // Next row id
		Prev   interface{} `json:"Prev,omitempty"`   // Prev row id
	}

	// PostRequest struct for mapping to post requests
	PostRequest struct {
		IO struct {
			Message string `json:"Message,omitempty"`
		} `json:"IO,omitempty"`
		Changes []map[string]interface{} `json:"Changes,omitempty"`
	}

	// Response struct for json responses
	PostResponse struct {
		IO struct {
			Message string
			Result  int
			Reason  string // Reason for error, for logical judgement
		}

		Changes []map[string]interface{}
		Body    []interface{}
	}
)

// ChangedErrorColor is the color for error row
const ChangedErrorColor = "rgb(255,0,0)"

// ChangedSuccessColor is the color for success row
const ChangedSuccessColor = "rgb(255, 255, 166)"

func GenColorChangeError(gr GridRow) ChangedRow {
	return ChangedRow{Id: gr.GetTreeGridID(), Color: ChangedErrorColor}
}

func GenMapColorChangeError(gr GridRow) map[string]interface{} {
	var inInterface map[string]interface{}
	change := GenColorChangeError(gr)
	inrec, _ := json.Marshal(change)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func CreateSuggestionResult(suggestionKey string, suggestion *Suggestion, tr *Treegrid) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	result[suggestionKey+"Suggest"] = suggestion
	result["id"] = tr.BodyParams.TreegridOriginID
	return result
}

// MakeResponseBody creates a new PostResponse with empty Changes and Body
func MakeResponseBody(resp *PostResponse) *PostResponse {
	if resp == nil {
		resp = &PostResponse{}
	}

	if resp.Changes == nil {
		resp.Changes = make([]map[string]interface{}, 0)
	}

	if resp.Body == nil {
		resp.Body = make([]interface{}, 0)
	}

	return resp
}

// GenGridRowChange get or generate a ChangedRow, which is used to return response for Upload API,
// generate a new error color ChangedRow if not exists,
// it will be set success if no error occurs in row processing(Add, Update, Delete)
func GenGridRowChange(gr GridRow) ChangedRow {
	if gr["ChangedRow"] != nil {
		return gr["ChangedRow"].(ChangedRow)
	}

	return ChangedRow{
		Id:      gr.GetTreeGridID(),
		NewId:   "",
		Changed: 0,
		Added:   0,
		Deleted: 0,
		Color:   ChangedErrorColor,
		Parent:  gr["Parent"],
		Next:    gr["Next"],
		Prev:    gr["Prev"],
	}
}

// SetGridRowChangedResult sets ChangedRow to GridRow
func SetGridRowChangedResult(gr GridRow, changedRow ChangedRow) {
	gr["ChangedRow"] = changedRow
}

// GetGridRowChangedResult gets ChangedRow map from GridRow
func GetGridRowChangedResult(gr GridRow) map[string]interface{} {
	changedRow, ok := gr["ChangedRow"].(ChangedRow)
	if !ok {
		return make(map[string]interface{})
	}

	return changedRow.ToMap()
}

// ToMap converts ChangedRow to map[string]interface{}
func (r ChangedRow) ToMap() map[string]interface{} {
	var m map[string]interface{}
	bytes, _ := json.Marshal(r)
	_ = json.Unmarshal(bytes, &m)
	return m
}
