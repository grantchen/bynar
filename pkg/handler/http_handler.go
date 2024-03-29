package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

var (
	// no permission error message
	noPermissionErrorMessage = "no permission"
	// no permission error reason
	noPermissionErrorReason = "NO_PERMISSION"
)

type HTTPTreeGridHandler struct {
	CallbackGetPageCountFunc CallBackGetPageCount
	CallbackGetPageDataFunc  CallBackGetPageData
	CallbackUploadDataFunc   CallBackUploadData
	CallBackGetCellDataFunc  CallBackGetCellData
	PathPrefix               string
	AccountManagerService    service.AccountManagerService
}

func (h *HTTPTreeGridHandler) getDB(_ *http.Request) *sql.DB {
	return nil
}

func (h *HTTPTreeGridHandler) HTTPHandleGetPageCount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	treegr, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	allPages, _ := h.CallbackGetPageCountFunc(treegr)

	response, err := json.Marshal(map[string]interface{}{
		"Body": []string{`#@@@` + fmt.Sprintf("%v", allPages)},
	})

	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	_, _ = w.Write(response)
}

func (h *HTTPTreeGridHandler) HTTPHandleGetPageData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	trGrid, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	var response = make([]map[string]string, 0, 100)

	response, _ = h.CallbackGetPageDataFunc(trGrid)

	var addData [][]map[string]string
	addData = append(addData, response)

	result, _ := json.Marshal(map[string][][]map[string]string{
		"Body": addData,
	})

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	_, _ = w.Write(result)
}

func (h *HTTPTreeGridHandler) HTTPHandleUpload(w http.ResponseWriter, r *http.Request) {
	var (
		postData = &treegrid.PostRequest{
			Changes: make([]map[string]interface{}, 10),
		}
	)

	// get and parse post data
	if err := r.ParseForm(); err != nil {
		logger.Debug("parse form err: ", err)
		writeErrorResponse(w, nil, err)
		return
	}

	if err := json.Unmarshal([]byte(r.Form.Get("Data")), &postData); err != nil {
		logger.Debug("unmarshal err: ", err)
		writeErrorResponse(w, nil, err)
		return
	}

	// b, _ := json.Marshal(postData)
	// logger.Debug("postData: ", string(b), "form data: ", r.Form.Get("Data"))
	resp, err := h.CallbackUploadDataFunc(postData)
	if err != nil {
		writeErrorResponse(w, resp, err)
		return
	}

	writeResponse(w, resp)
}

func (h *HTTPTreeGridHandler) HTTPHandleCell(w http.ResponseWriter, r *http.Request) {
	logger.Debug("request come here")
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	trGrid, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	resp, err := h.CallBackGetCellDataFunc(trGrid)
	if err != nil {
		writeErrorResponse(w, resp, err)
		return
	}

	writeResponse(w, resp)
}

// response no permission error
func writeNoPermissionResponse(w http.ResponseWriter, resp *treegrid.PostResponse) {
	if resp == nil {
		resp = &treegrid.PostResponse{}
	}

	resp.IO.Result = -1
	resp.IO.Message = noPermissionErrorMessage
	resp.IO.Reason = noPermissionErrorReason

	// write response with error
	writeResponse(w, resp)
}

// WriteErrorResponse
// if err occures then maybe invalid request so:
// * log error
// resp.IO.Result = -1 - need for treegrid to mark as error (negative numbers are considered as errors)
// resp.IO.Message message for treegrid for modal window, for production better not use err message
func writeErrorResponse(w http.ResponseWriter, resp *treegrid.PostResponse, err error) {
	if resp == nil {
		resp = &treegrid.PostResponse{}
	}

	if err != nil {
		log.Println("Err", err)

		resp.IO.Result = -1
		resp.IO.Message = err.Error()
	}

	// write response with error
	writeResponse(w, resp)
}

// WriteResponse writes json response
func writeResponse(w http.ResponseWriter, resp *treegrid.PostResponse) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", " application/json")
	w.WriteHeader(http.StatusOK)

	resp = treegrid.MakeResponseBody(resp)

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
