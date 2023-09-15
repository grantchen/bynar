package render

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// CorsMiddleware solve the CORS problem
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			Ok(w, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// DecodeJSON decode and validate the json request
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

// Ok success response
func Ok(w http.ResponseWriter, v interface{}) {
	if v == nil {
		v = struct {
			Msg string `json:"msg"`
		}{Msg: "OK"}
	}
	renderJSON(w, v)
}

// Error failure response
func Error(w http.ResponseWriter, msg string) {
	data, err := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: msg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	http.Error(w, string(data), http.StatusInternalServerError)
}

func MethodNotAllowed(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Write(data)
}
