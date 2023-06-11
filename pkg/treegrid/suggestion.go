package treegrid

type Suggestion struct {
	Items []map[string]interface{}
}

type SuggestionItemRow struct {
	Columns int
	Item    SuggestionItem
	Value   string
}

type SuggestionItem struct {
	Name  string
	Value string
}

// func CreateSuggestion(mapData []map[string]interface{}, column []string) Suggestion {
// 	header
// }
