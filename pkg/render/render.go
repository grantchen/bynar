package render

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

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

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
