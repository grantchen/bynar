package http_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

type TransferHttpHandler struct {
	transferService service.TransferService
}

func NewTransferHttpHandler(transferService service.TransferService) *TransferHttpHandler {
	return &TransferHttpHandler{transferService: transferService}
}

func (t *TransferHttpHandler) HandleGetPageCount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	trRequest, err := treegrid_model.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Fatal(err)
	}

	treegr, err := treegrid_model.NewTreegrid(trRequest)
	if err != nil {
		log.Fatal(err)
	}

	allPages := t.transferService.GetPagesCount(treegr)

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

func (t *TransferHttpHandler) HandleGetPageTransferData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	trRequest, err := treegrid_model.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Fatal(err)
	}

	trGrid, err := treegrid_model.NewTreegrid(trRequest)
	if err != nil {
		log.Fatal(err)
	}

	var response = make([]map[string]string, 0, 100)

	response, _ = t.transferService.GetTransfersPageData(trGrid)

	addData := [][]map[string]string{}
	addData = append(addData, response)

	result, _ := json.Marshal(map[string][][]map[string]string{
		"Body": addData,
	})

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.Write(result)
}
