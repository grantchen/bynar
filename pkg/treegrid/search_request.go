package treegrid

import "encoding/json"

// Request represents request made by treegrid
type Request struct {
	Cfg     RequestConfig            `json:"Cfg,omitempty"`
	Filters []map[string]interface{} `json:"Filters,omitempty"`
	Body    []BodyParam              `json:"Body,omitempty"`
}

// ParseRequest parses request from input bytes
func ParseRequest(jsonBuf []byte) (r *Request, err error) {
	err = json.Unmarshal(jsonBuf, &r)

	return
}

// RequestConfig contains all request's query params
type RequestConfig struct {
	Sort string `json:"Sort,omitempty"`
	// Data example: "col1,col2"
	SortCols            string `json:"SortCols,omitempty"`
	SortTypes           string `json:"SortTypes,omitempty"`
	Group               string `json:"Group,omitempty"`
	GroupCols           string `json:"GroupCols,omitempty"`
	SearchAction        string `json:"SearchAction,omitempty"`
	SearchExpression    string `json:"SearchExpression,omitempty"`
	SearchType          string `json:"SearchType,omitempty"`
	SearchCaseSensitive string `json:"SearchCase_Sensitive,omitempty"`
	SearchCells         string `json:"SearchCells,omitempty"`
	SearchMethod        string `json:"SearchMethod,omitempty"`
	SearchDefs          string `json:"SearchDefs,omitempty"`
	SearchCols          string `json:"SearchCols,omitempty"`
	TimeZone            string `json:"TimeZone,omitempty"`
}
