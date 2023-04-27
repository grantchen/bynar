package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type CallBackGetPageCount func(tr *treegrid.Treegrid) float64
type CallBackGetPageData func(tr *treegrid.Treegrid) ([]map[string]string, error)
type CallBackUploadData func(req *treegrid.PostRequest) (*treegrid.PostResponse, error)
type CallBackGetCellData func(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error)

type HTTPTreeGridHandler struct {
	CallbackGetPageCountFunc CallBackGetPageCount
	CallbackGetPageDataFunc  CallBackGetPageData
	CallbackUploadDataFunc   CallBackUploadData
	CallBackGetCellDataFunc  CallBackGetCellData
}

func (h *HTTPTreeGridHandler) HTTPHandleGetPageCount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Fatal(err)
	}

	treegr, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Fatal(err)
	}

	allPages := h.CallbackGetPageCountFunc(treegr)

	response, err := json.Marshal((map[string]interface{}{
		"Body": []string{`#@@@` + fmt.Sprintf("%v", allPages)},
	}))

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func (h *HTTPTreeGridHandler) HTTPHandleGetPageTransferData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Fatal(err)
	}

	trGrid, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Fatal(err)
	}

	var response = make([]map[string]string, 0, 100)

	response, _ = h.CallbackGetPageDataFunc(trGrid)

	addData := [][]map[string]string{}
	addData = append(addData, response)

	result, _ := json.Marshal(map[string][][]map[string]string{
		"Body": addData,
	})

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.Write(result)
}

func (h *HTTPTreeGridHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	var (
		postData = &treegrid.PostRequest{
			Changes: make([]map[string]interface{}, 10),
		}

		resp = &treegrid.PostResponse{}
	)

	// get and parse post data
	if err := r.ParseForm(); err != nil {
		writeErrorResponse(w, resp, err)

		return
	}

	if err := json.Unmarshal([]byte(r.Form.Get("Data")), &postData); err != nil {
		writeErrorResponse(w, resp, err)

		return
	}

	resp, err := h.CallbackUploadDataFunc(postData)
	if err != nil {
		writeErrorResponse(w, resp, err)

		return
	}

	writeResponse(w, resp)
}

func (h *HTTPTreeGridHandler) HandleCell(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Fatal(err)
	}

	trGrid, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := h.CallBackGetCellDataFunc(r.Context(), trGrid)
	if err != nil {
		writeErrorResponse(w, resp, err)

		return
	}

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
