/**
    @author: dongjs
    @date: 2023/9/15
    @description:
**/

package handler

import (
	"context"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type HTTPHandlerWithDynamicDB struct {
	Path           string
	ConnectionPool ConnectionResolver
	RequestFunc    func(w http.ResponseWriter, r *http.Request)
}

func (h *HTTPHandlerWithDynamicDB) HandleHTTPReqWithDynamicDB() {

	http.Handle(h.Path, render.CorsMiddleware(h.verifyIdTokenAndInitDynamicDB(http.HandlerFunc(h.RequestFunc))))

}

func (h *HTTPHandlerWithDynamicDB) verifyIdTokenAndInitDynamicDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defaultResponse := &treegrid.PostResponse{}
		defaultResponse.Changes = make([]map[string]interface{}, 0)

		code, msg, claims := middleware.VerifyIdToken(r)
		if http.StatusOK != code {
			if "" == msg {
				msg = http.StatusText(code)
			}
			http.Error(w, msg, code)
			return
		}

		connString, err := h.getDynamicDBConnection(claims.TenantUuid, claims.OrganizationUuid)
		if err != nil {
			logrus.Errorf("verifyIdTokenAndInitDynamicDB: getDynamicDBConnection error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db, err := h.ConnectionPool.Get(connString)

		if err != nil {
			logrus.Errorf("verifyIdTokenAndInitDynamicDB: get connection db error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqContext := &middleware.TokenAndDyDynamicDBContext{
			ConnectionString: connString,
			DynamicDB:        db,
			Claims:           claims,
		}
		ctx := context.WithValue(r.Context(), middleware.IdTokenAndDynamicDBRequestContextKey, reqContext)
		newReq := r.WithContext(ctx)
		next.ServeHTTP(w, newReq)
	})
}

func (h *HTTPHandlerWithDynamicDB) getDynamicDBConnection(tenantUuid, organizationUuid string) (string, error) {
	if len(os.Getenv(tenantUuid)) == 0 {
		return "", errors.New("no mysql conn environment of " + tenantUuid)
	}
	envs := strings.Split(os.Getenv(tenantUuid), "/")
	connStr := envs[0] + "/" + organizationUuid
	if len(envs) > 1 {
		connStr += envs[1]
	}
	return connStr, nil
}
