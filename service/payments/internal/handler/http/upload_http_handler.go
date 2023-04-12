package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload/factory"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/treegrid"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
)

type UploadHandler struct {
	ModuleID  int
	AccountID int
}

func (u *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	var (
		postData = &treegrid.PostRequest{
			Changes: make([]map[string]interface{}, 10),
		}

		resp = &treegrid.PostResponse{}
	)

	// get and parse post data
	if err := r.ParseForm(); err != nil {
		WriteErrorResponse(w, resp, err)

		return
	}

	if err := json.Unmarshal([]byte(r.Form.Get("Data")), &postData); err != nil {
		WriteErrorResponse(w, resp, err)

		return
	}

	// TODO: Handle in the similar way as AWS Lambda, secret manager and connection pool might be used.
	// We're not touching this code because it isn't a part of specification.
	conn := sql_db.Conn()

	uploadSvc, err := factory.NewUploadService(conn, u.ModuleID, u.AccountID)
	if err != nil {
		log.Println(err)

		WriteErrorResponse(w, resp, errors.New("problem while initializing uploader service"))

		return
	}

	resp, err = uploadSvc.Handle(postData)
	if err != nil {
		WriteErrorResponse(w, resp, err)

		return
	}

	WriteResponse(w, resp)
}

func (u *UploadHandler) HandleUploadTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-type", " application/json")
	var postData = &treegrid.PostRequest{
		Changes: make([]map[string]interface{}, 10),
	}
	var resp = &treegrid.PostResponse{}
	if err := r.ParseForm(); err != nil {

		return
	}
	fmt.Printf("method: %s\n", r.Method)
	fmt.Printf("data: %s\n", r.Form.Get("Data"))
	json.Unmarshal([]byte(r.Form.Get("Data")), &postData)

	WriteErrorResponse(w, resp, nil)
}

// WriteErrorResponse
// if err occures then maybe invalid request so:
// * log error
// resp.IO.Result = -1 - need for treegrid to mark as error (negative numbers are considered as errors)
// resp.IO.Message message for treegrid for modal window, for production better not use err message
func WriteErrorResponse(w http.ResponseWriter, resp *treegrid.PostResponse, err error) {
	if resp == nil {
		resp = &treegrid.PostResponse{}
	}

	if err != nil {
		log.Println("Err", err)

		resp.IO.Result = -1
		resp.IO.Message = err.Error()
	}

	// write response with error
	WriteResponse(w, resp)
}

// WriteResponse writes json response
func WriteResponse(w http.ResponseWriter, resp *treegrid.PostResponse) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", " application/json")
	w.WriteHeader(http.StatusOK)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("Err", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Println("Err", err)
	}
}
