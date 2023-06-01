package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type HTTPTreeGridHandlerWithDynamicDB struct {
	PathPrefix             string
	AccountManagerService  service.AccountManagerService
	TreeGridServiceFactory treegrid.TreeGridServiceFactoryFunc
}

type ReqContext struct {
	connectionString string
	db               *sql.DB
}

type key string

const RequestContextKey key = "reqContext"

func (h *HTTPTreeGridHandlerWithDynamicDB) getDB(r *http.Request) *sql.DB {
	reqContext := r.Context().Value(RequestContextKey).(*ReqContext)
	return reqContext.db
}

func (h *HTTPTreeGridHandlerWithDynamicDB) getTreeGridService(r *http.Request) treegrid.TreeGridService {
	db := h.getDB(r)
	return h.TreeGridServiceFactory(db)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HTTPHandleGetPageCount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Println(err)
	}

	treegr, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Println(err)
	}

	treegridService := h.getTreeGridService(r)
	allPages := treegridService.GetPageCount(treegr)

	response, err := json.Marshal((map[string]interface{}{
		"Body": []string{`#@@@` + fmt.Sprintf("%v", allPages)},
	}))

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HTTPHandleGetPageData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Println(err)
	}

	trGrid, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Println(err)
	}

	var response = make([]map[string]string, 0, 100)

	treegridService := h.getTreeGridService(r)
	response, _ = treegridService.GetPageData(trGrid)

	addData := [][]map[string]string{}
	addData = append(addData, response)

	result, _ := json.Marshal(map[string][][]map[string]string{
		"Body": addData,
	})

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.Write(result)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HTTPHandleUpload(w http.ResponseWriter, r *http.Request) {
	var (
		postData = &treegrid.PostRequest{
			Changes: make([]map[string]interface{}, 10),
		}

		resp = &treegrid.PostResponse{}
	)

	// get and parse post data
	if err := r.ParseForm(); err != nil {
		logger.Debug("parse form err: ", err)
		writeErrorResponse(w, resp, err)

		return
	}

	if err := json.Unmarshal([]byte(r.Form.Get("Data")), &postData); err != nil {
		logger.Debug("unmarshal err: ", err)
		writeErrorResponse(w, resp, err)

		return
	}

	treegridService := h.getTreeGridService(r)
	resp, err := treegridService.Upload(postData)

	if err != nil {
		writeErrorResponse(w, resp, err)

		return
	}

	writeResponse(w, resp)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HandleCell(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

	trRequest, err := treegrid.ParseRequest([]byte(r.Form.Get("Data")))
	if err != nil {
		log.Println(err)
	}

	trGrid, err := treegrid.NewTreegrid(trRequest)
	if err != nil {
		log.Println(err)
	}

	treegridService := h.getTreeGridService(r)
	resp, err := treegridService.GetCellData(r.Context(), trGrid)
	if err != nil {
		writeErrorResponse(w, resp, err)

		return
	}

	writeResponse(w, resp)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) authenMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := "<<Extract token from req here>>"

		// if true {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		logger.Debug("check permission")
		permission, ok, err := h.AccountManagerService.CheckPermission(token)

		if err != nil {
			log.Println("Err", err)
			writeErrorResponse(w, &treegrid.PostResponse{}, err)
			return
		}

		if !ok {
			writeErrorResponse(w, &treegrid.PostResponse{}, err)
			return
		}
		var connString string
		connString, _ = h.AccountManagerService.GetNewStringConnection(token, permission)
		ctx := context.WithValue(r.Context(), RequestContextKey, &ReqContext{connectionString: connString})
		newReq := r.WithContext(ctx)
		next.ServeHTTP(w, newReq)
	})
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HandleHTTPReqWithAuthenMWAndDefaultPath() {

	if h.AccountManagerService == nil {
		panic("account manager service is null")
	}

	http.Handle(h.PathPrefix+"/upload", h.authenMW(http.HandlerFunc(h.HTTPHandleUpload)))
	http.Handle(h.PathPrefix+"/data", h.authenMW(http.HandlerFunc(h.HTTPHandleGetPageCount)))
	http.Handle(h.PathPrefix+"/page", h.authenMW(http.HandlerFunc(h.HTTPHandleGetPageData)))

}
