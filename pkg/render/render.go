package render

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
)

type PaginationResponse struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
	PerPage     int `json:"per_page"`
}

func newPagination(currentPage, totalCount, perPage int) *PaginationResponse {
	return &PaginationResponse{
		CurrentPage: currentPage,
		TotalCount:  totalCount,
		PerPage:     perPage,
		TotalPages:  int(math.Ceil(float64(totalCount) / float64(perPage))),
	}
}

type Response struct {
	Code       int                 `json:"code"`
	Message    string              `json:"message"`
	Data       interface{}         `json:"data,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

func DecodeJSON(r io.Reader, v interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func Ok(w http.ResponseWriter, v interface{}) {
	renderJSON(w, Response{Code: 200, Message: "OK", Data: v})
}

func Pagination(w http.ResponseWriter, v []interface{}, currentPage, totalCount, perPage int) {
	renderJSON(w, Response{Code: 200, Message: "OK", Data: v, Pagination: newPagination(currentPage, totalCount, perPage)})
}

func Error(w http.ResponseWriter, code int, msg string) {
	renderJSON(w, Response{Code: code, Message: msg})
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
