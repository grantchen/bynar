/**
    @author: dongjs
    @date: 2023/9/15
    @description:
**/

package handler

import (
	"context"
	"errors"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

// VerifyIdTokenAndInitDynamicDB verify idToken correct and create dynamic db connection
func VerifyIdTokenAndInitDynamicDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, msg, claims := middleware.VerifyIdToken(r)
		if http.StatusOK != code {
			if "" == msg {
				msg = http.StatusText(code)
			}
			http.Error(w, msg, code)
			return
		}

		connString, err := getDynamicDBConnection(claims.TenantUuid, claims.OrganizationUuid)
		if err != nil {
			logrus.Errorf("VerifyIdTokenAndInitDynamicDB: getDynamicDBConnection error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db, err := sql_db.NewConnection(connString)

		if err != nil {
			logrus.Errorf("VerifyIdTokenAndInitDynamicDB: get connection db error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
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

// get dynamic db connection url
func getDynamicDBConnection(tenantUuid, organizationUuid string) (string, error) {
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
