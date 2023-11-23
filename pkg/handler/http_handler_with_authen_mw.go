package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type ConnectionResolver interface {
	Get(string) (*sql.DB, error)
}

const UploadPathString = "upload"
const PageCountPathString = "data"
const PageDataPathString = "page"
const CellDataPathString = "cell"

type HTTPTreeGridHandlerWithDynamicDB struct {
	PathPrefix             string
	AccountManagerService  service.AccountManagerService
	TreeGridServiceFactory treegrid.TreeGridServiceFactoryFunc
	ConnectionPool         ConnectionResolver
	IsValidatePermissions  bool
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
	OrganizationUuid string
	language         string
	//claims           *middleware.IdTokenClaims
}

type ModulePath struct {
	module      string // transfers or payments
	pathFeature string // upload, data, cell
}
type key string

const RequestContextKey key = "reqContext"

var PolicyMap = map[string][]string{
	"list":   {"data", "page", "upload", "cell"},
	"add":    {"upload:Added", "upload"},
	"update": {"upload:Changed", "upload"},
	"delete": {"upload:Deleted", "upload"},
}

func (h *HTTPTreeGridHandlerWithDynamicDB) getRequestContext(r *http.Request) *ReqContext {
	reqContext := r.Context().Value(RequestContextKey).(*ReqContext)
	return reqContext
}

func (h *HTTPTreeGridHandlerWithDynamicDB) getTreeGridService(r *http.Request) treegrid.TreeGridService {
	reqContext := h.getRequestContext(r)
	return h.TreeGridServiceFactory(reqContext.db, reqContext.AccountID, reqContext.OrganizationUuid, reqContext.PermissionInfo, reqContext.language)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HTTPHandleGetPageCount(w http.ResponseWriter, r *http.Request) {
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

	treegridService := h.getTreeGridService(r)
	allPages, err := treegridService.GetPageCount(treegr)

	if err != nil {
		writeErrorResponse(w, nil, err)
		return
	}

	response, err := json.Marshal((map[string]interface{}{
		"Body": []string{`#@@@` + fmt.Sprintf("%v", allPages)},
	}))

	if err != nil {
		log.Println(err)
		writeErrorResponse(w, nil, err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HTTPHandleGetPageData(w http.ResponseWriter, r *http.Request) {
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

	treegridService := h.getTreeGridService(r)
	response, err = treegridService.GetPageData(trGrid)
	if err != nil {
		writeErrorResponse(w, nil, err)
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

	treegridService := h.getTreeGridService(r)
	resp, err := treegridService.Upload(postData)

	if err != nil {
		writeErrorResponse(w, resp, err)
		return
	}

	writeResponse(w, resp)
}

func (h *HTTPTreeGridHandlerWithDynamicDB) HTTPHandleCell(w http.ResponseWriter, r *http.Request) {
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
	data, _ := io.ReadAll(r.Body)
	query, _ := url.QueryUnescape(string(data))
	logrus.Info("module data ", query)
	if modulePath.pathFeature == "upload" {
		if strings.Contains(query, "Group") {
			modulePath.pathFeature = "data"
		} else if strings.Contains(query, "Added") {
			modulePath.pathFeature += ":Added"
		} else if strings.Contains(query, "Deleted") {
			modulePath.pathFeature += ":Deleted"
		} else if strings.Contains(query, "Changed") {
			modulePath.pathFeature += ":Changed"
		}
	}
	r.Body.Close() //  must close
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	return modulePath
}

func (h *HTTPTreeGridHandlerWithDynamicDB) authenMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defaultResponse := &treegrid.PostResponse{}
		defaultResponse.Changes = make([]map[string]interface{}, 0)
		defaultResponse.Body = make([]interface{}, 0)

		code, claims, err := middleware.VerifyIdToken(r)
		if http.StatusOK != code {
			writeErrorResponse(w, defaultResponse, errors.New(http.StatusText(code)))
			return
		}
		if !claims.OrganizationStatus || !claims.TenantStatus || claims.TenantSuspended {
			writeErrorResponse(w, defaultResponse, errors.New("no permission"))
			return
		}
		modulePath := getModuleFromPath(r)
		if !claims.OrganizationAccount && modulePath.module == "invoices" {
			writeErrorResponse(w, defaultResponse, errors.New("no permission"))
			return
		}
		var connString string
		// Initialize to the "accounts_management" database connection
		db := sql_db.Conn()
		// Validate permissions
		if h.IsValidatePermissions {
			logger.Debug("check permission")
			permission := &repository.PermissionInfo{}
			connString, _ = h.AccountManagerService.GetNewStringConnection(claims.TenantUuid, claims.OrganizationUuid, permission)
			db, err = h.ConnectionPool.Get(connString)
			if err != nil {
				log.Println("Err get policy", err)
				writeErrorResponse(w, defaultResponse, errors.New("no permission"))
				return
			}
			var val string
			err = db.QueryRow("SELECT policies FROM users WHERE users.email = ?", claims.Email).Scan(&val)
			if err != nil {
				log.Println("Err get policy", err)
				writeErrorResponse(w, defaultResponse, errors.New("no permission"))
				return
			}
			policy := models.Policy{Services: make([]models.ServicePolicy, 0)}
			err = json.Unmarshal([]byte(val), &policy)
			if err != nil {
				log.Println("Err get policy", err)
				writeErrorResponse(w, defaultResponse, errors.New("no permission"))
				return
			}
			allowed := false
			for _, i := range policy.Services {
				if i.Name == "*" || i.Name == modulePath.module {
					for _, j := range i.Permissions {
						if j == "*" {
							allowed = true
							break
						} else {
							val, _ := PolicyMap[j]
							for _, m := range val {
								if m == modulePath.pathFeature {
									allowed = true
									break
								}
							}
						}
					}
				}
			}
			if !allowed {
				log.Println("not allowed to get policy "+modulePath.module, modulePath.pathFeature)
				writeErrorResponse(w, defaultResponse, errors.New("no permission"))
				return
			}
		}

		reqContext := &ReqContext{
			connectionString: connString,
			db:               db,
			AccountID:        claims.AccountId,
			PermissionInfo: &treegrid.PermissionInfo{
				IsAccessAll: true,
			},
			OrganizationUuid: claims.OrganizationUuid,
			language:         claims.Language,
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
	http.Handle(h.PathPrefix+"/"+CellDataPathString, render.CorsMiddleware(h.authenMW(http.HandlerFunc(h.HTTPHandleCell))))

}
