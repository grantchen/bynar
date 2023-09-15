/**
    @author: dongjs
    @date: 2023/9/15
    @description:
**/

package handler

import (
	"context"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"net/http"
)

// VerifyIdToken Verify if the token is correct
func VerifyIdToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify if the token is correct
		code, msg, claims := middleware.VerifyIdToken(r)
		if http.StatusOK != code {
			if "" == msg {
				msg = http.StatusText(code)
			}
			http.Error(w, msg, code)
			return
		}
		ctx := r.Context()
		if claims != nil {
			ctx = context.WithValue(ctx, "id_token", *claims)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
