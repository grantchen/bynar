/**
    @author: dongjs
    @date: 2023/9/15
    @description:
**/

package handler

import (
	"context"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"net/http"
	"os"
	"strings"
)

// VerifyIdTokenAndInitDynamicDB verify idToken correct and create dynamic db connection
func VerifyIdTokenAndInitDynamicDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, claims, err := middleware.VerifyIdToken(r)
		if err != nil {
			LogInternalError(errors.NewUnknownError("verify id_token fail").WithInternalCause(err))
		}
		if http.StatusOK != code {
			http.Error(w, http.StatusText(code), code)
			return
		}

		connString, err := getDynamicDBConnection(claims.TenantUuid, claims.OrganizationUuid)
		if err != nil {
			LogInternalError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db, err := sql_db.InitializeConnection(connString)

		if err != nil {
			LogInternalError(errors.NewUnknownError("new dynamic db connection fail").WithInternalCause(err))
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
		return "", errors.NewUnknownError("no mysql conn environment of " + tenantUuid).WithInternal()
	}
	envs := strings.Split(os.Getenv(tenantUuid), "/")
	connStr := envs[0] + "/" + organizationUuid
	if len(envs) > 1 {
		connStr += envs[1]
	}
	return connStr, nil
}
