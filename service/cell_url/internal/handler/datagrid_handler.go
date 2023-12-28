package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/service"
)

type HTTPDataGridHandler struct {
	dataGridService service.DataGridService
}

func NewHTTPDataGridHandler(dataGridService service.DataGridService) *HTTPDataGridHandler {
	return &HTTPDataGridHandler{
		dataGridService: dataGridService,
	}
}

func (h *HTTPDataGridHandler) GetTreeGridViewData(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	context := r.Context()

	// Check if the "Data" form value is not empty
	if r.FormValue("Data") != "" {
		// Extract the "Data" form value and remove any trailing comma in the JSON string
		tgCell := r.FormValue("Data")
		tgCell = strings.Replace(tgCell, "Body\":[{,", "Body\":[{", 1)

		// Parse the JSON string and extract the "Val" and "id" fields
		m := make(map[string]interface{})
		_ = json.Unmarshal([]byte(tgCell), &m)
		body := m["Body"].([]interface{})
		mapParam := body[0].(map[string]interface{})
		val := mapParam["Val"].(string)
		id := mapParam["id"].(string)

		// Call the GetGridDataByValue function to get grid data for the search value
		gridData, err := h.dataGridService.GetGridDataByValue(context, val, id)
		if err != nil {
			// Return 500 Internal Server Error if there is an error during retrieval
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate a JSON response and write the encoded grid data to the response writer
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(gridData)
	}
}
