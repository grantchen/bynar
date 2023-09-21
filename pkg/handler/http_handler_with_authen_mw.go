package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type ConnectionResolver interface {
	Get(string) (*sql.DB, error)
}

const UploadPathString = "upload"
const PageCountPathString = "data"
const PageDataPathString = "page"

type HTTPTreeGridHandlerWithDynamicDB struct {
	PathPrefix             string
	AccountManagerService  service.AccountManagerService
	TreeGridServiceFactory treegrid.TreeGridServiceFactoryFunc
	ConnectionPool         ConnectionResolver
}

func NewHTTPTreeGridHandlerWithDynamicDB(
	accountManagerService service.AccountManagerService,
	treeGridServiceFactory treegrid.TreeGridServiceFactoryFunc,
	connectionPool ConnectionResolver,
) *HTTPTreeGridHandlerWithDynamicDB {
	return &HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountManagerService,
		TreeGridServiceFactory: treeGridServiceFactory,
		ConnectionPool:         connectionPool,
	}
}

type ReqContext struct {
	connectionString string
	db               *sql.DB
	AccountID        int
	PermissionInfo   *treegrid.PermissionInfo
}

type ModulePath struct {
	module      string // transfers or payments
	pathFeature string // upload, data, cell
}
type key string

const RequestContextKey key = "reqContext"

func (h *HTTPTreeGridHandlerWithDynamicDB) getRequestContext(r *http.Request) *ReqContext {
	reqContext := r.Context().Value(RequestContextKey).(*ReqContext)
	return reqContext
}

func (h *HTTPTreeGridHandlerWithDynamicDB) getTreeGridService(r *http.Request) treegrid.TreeGridService {
	reqContext := h.getRequestContext(r)
	return h.TreeGridServiceFactory(reqContext.db, reqContext.AccountID, reqContext.PermissionInfo)
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
	allPages, err := treegridService.GetPageCount(treegr)

	if err != nil {
		defaultResponse := &treegrid.PostResponse{}
		defaultResponse.Changes = make([]map[string]interface{}, 0)
		writeErrorResponse(w, defaultResponse, err)
		return
	}

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
	response, err = treegridService.GetPageData(trGrid)
	if err != nil {
		defaultResponse := &treegrid.PostResponse{}
		defaultResponse.Changes = make([]map[string]interface{}, 0)
		writeErrorResponse(w, defaultResponse, err)
		return
	}

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

func getModuleFromPath(r *http.Request) *ModulePath {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")
	modulePath := &ModulePath{}
	modulePath.pathFeature = splittedPath[len(splittedPath)-1]
	if len(splittedPath) > 1 {
		modulePath.module = splittedPath[len(splittedPath)-2]
	} else {
		modulePath.module = splittedPath[0]
	}
	return modulePath
}

func (h *HTTPTreeGridHandlerWithDynamicDB) authenMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defaultResponse := &treegrid.PostResponse{}
		defaultResponse.Changes = make([]map[string]interface{}, 0)

		code, claims, err := middleware.VerifyIdToken(r)
		if http.StatusOK != code {
			writeErrorResponse(w, defaultResponse, errors.New(http.StatusText(code)))
			return
		}
		if !claims.OrganizationStatus || !claims.TenantStatus || claims.TenantSuspended {
			writeErrorResponse(w, defaultResponse, errors.New("no permission"))
			return
		}

		logger.Debug("check permission")
		permission := &repository.PermissionInfo{}
		// TODO:
		// permission, ok, err := h.AccountManagerService.CheckPermission(claims)

		// if err != nil {
		// 	log.Println("Err", err)
		// 	writeErrorResponse(w, defaultResponse, err)
		// 	return
		// }

		// if !ok {
		// 	writeErrorResponse(w, defaultResponse, err)
		// 	return
		// }

		// check role TODO:
		// roles, err := h.AccountManagerService.GetRole(0)

		// if err != nil {
		// 	log.Println("Err", err)
		// 	writeErrorResponse(w, defaultResponse, err)
		// 	return
		// }

		// logger.Debug("role: ", roles, "req string: ", r.URL.Path, "module str: ", modulePath.pathFeature)

		// moduleVal, ok := roles[modulePath.module]
		// if !ok {
		// 	writeErrorResponse(w, defaultResponse, fmt.Errorf("not found module in policies: [%s]", modulePath))
		// 	return
		// }

		// use for pass to modules to filter permission, 0 mean have all permission
		// accID := 0
		// if moduleVal == 0 {
		// 	writeErrorResponse(w, defaultResponse, fmt.Errorf("no permission allowed to access module: [%s]", modulePath.module))
		// 	return
		// }

		// moduleDataVal, ok := roles[modulePath.module+"_data"]
		// if !ok {
		// 	writeErrorResponse(w, defaultResponse, fmt.Errorf("not found module data in policies: [%s]", modulePath.module+"_data"))
		// 	return
		// }
		// accID = 0

		// user can access all module
		// if moduleDataVal == 1 {
		// 	accID = 0
		// } else {
		// 	if modulePath.pathFeature != PageCountPathString && modulePath.pathFeature != PageDataPathString {
		// 		writeErrorResponse(w, defaultResponse, fmt.Errorf("action is not allowed, Only /page and /data allowed"))
		// 		return
		// 	}
		// }

		var connString string

		connString, _ = h.AccountManagerService.GetNewStringConnection(claims.TenantUuid, claims.OrganizationUuid, permission)
		//hardcode to test
		// connString = "root:123456@tcp(localhost:3306)/172c1ecd-fd74-40a2-8717-1cc9a5e18880"

		db, err := h.ConnectionPool.Get(connString)

		if err != nil {
			log.Println("Err get connection db", err)
			writeErrorResponse(w, defaultResponse, err)
			return
		}
		var status int
		db.QueryRow(`SELECT status FROM users WHERE email = ?`, claims.Email).Scan(&status)
		if status == 0 {
			writeErrorResponse(w, defaultResponse, errors.New("user is disabled"))
			return
		}

		modulePath := getModuleFromPath(r)
		var val int
		// TODO: no policy in db at current time
		err = db.QueryRow(fmt.Sprintf("SELECT policies.%s FROM users LEFT JOIN policies ON policies.id = users.policy_id WHERE users.email = ?", modulePath.module), claims.Email).Scan(&val)
		if err != nil {
			log.Println("Err get policy", err)
			// writeErrorResponse(w, defaultResponse, err)
			// return
		}
		if val != 1 {
			log.Println("not allowed to get policy " + modulePath.module)
			// writeErrorResponse(w, defaultResponse, errors.New("do not have policy"))
			// return
		}

		reqContext := &ReqContext{
			connectionString: connString,
			db:               db,
			AccountID:        claims.OrganizationUserId,
			PermissionInfo: &treegrid.PermissionInfo{
				IsAccessAll: true,
			},
		}
		ctx := context.WithValue(r.Context(), RequestContextKey, reqContext)
		newReq := r.WithContext(ctx)
		next.ServeHTTP(w, newReq)
	})
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HandleHTTPReqWithAuthenMWAndDefaultPath() {

	if h.AccountManagerService == nil {
		panic("account manager service is null")
	}

	logger.Debug(h.PathPrefix + "/" + UploadPathString)
	http.Handle(h.PathPrefix+"/"+UploadPathString, render.CorsMiddleware(h.authenMW(http.HandlerFunc(h.HTTPHandleUpload))))
	http.Handle(h.PathPrefix+"/"+PageCountPathString, render.CorsMiddleware(h.authenMW(http.HandlerFunc(h.HTTPHandleGetPageCount))))
	http.Handle(h.PathPrefix+"/"+PageDataPathString, render.CorsMiddleware(h.authenMW(http.HandlerFunc(h.HTTPHandleGetPageData))))

}
