package treegrid

type (
	// ChangedRow: used to return Messages for POST update
	ChangedRow struct {
		Id      string `json:"id,omitempty"`
		NewId   string `json:"NewId,omitempty"`
		Changed int    `json:"Changed,omitempty"`
		Added   int    `json:"Added,omitempty"`
		Deleted int    `json:"Deleted,omitempty"`
		Color   string `json:"Color,omitempty"`
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
	}
)

type ChangeItemType string

const (
	ChangeTypeNode ChangeItemType = "Node"
	ChangeTypeItem ChangeItemType = "Data"
)
