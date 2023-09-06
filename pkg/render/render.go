package render

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
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
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	typ := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		if _, ok := typ.Field(i).Tag.Lookup("valid"); ok {
			if len(val.Field(i).String()) == 0 {
				return fmt.Errorf("%s is empty", typ.Field(i).Name)
			}
		}
	}
	return err
}

func Ok(w http.ResponseWriter, v interface{}) {
	if v == nil {
		v = struct {
			Msg string `json:"msg"`
		}{Msg: "OK"}
	}
	renderJSON(w, v)
}

func Error(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

func Pagination(w http.ResponseWriter, v []interface{}, currentPage, totalCount, perPage int) {
	renderJSON(w, Response{Code: 200, Message: "OK", Data: v, Pagination: newPagination(currentPage, totalCount, perPage)})
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
