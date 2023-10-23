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
		}

		Changes []map[string]interface{}
		Body    []interface{}
	}
)

func GenColorChangeError(gr GridRow) ChangedRow {
	return ChangedRow{Id: gr.GetTreeGridID(), Color: "rgb(255,0,0)"}
}

func GenColorChangeSuccess(gr GridRow) ChangedRow {
	return ChangedRow{Id: gr.GetTreeGridID(), Color: "rgb(255, 255, 166)"}
}

func GenMapColorChangeError(gr GridRow) map[string]interface{} {
	var inInterface map[string]interface{}
	change := GenColorChangeError(gr)
	inrec, _ := json.Marshal(change)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func GenMapColorChangeSuccess(gr GridRow) map[string]interface{} {
	var inInterface map[string]interface{}
	change := GenColorChangeSuccess(gr)
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

type ChangeItemType string

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
