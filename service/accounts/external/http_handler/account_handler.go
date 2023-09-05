package http_handler

import (
	"database/sql"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
)

type AccountHandler struct {
	ac service.AccountService
}

func NewHTTPHandler(db *sql.DB) *AccountHandler {
	ac := service.NewUserService(db)
	return &AccountHandler{ac}
}

func (h *AccountHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		var req model.SignupRequest
		// TODO:
		render.DecodeJSON(r.Body, req)
		h.ac.CreateUser(req.Email)
		render.Ok(w, "OK")
	}
}

func (h *AccountHandler) Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		// TODO:
		render.Ok(w, "OK")
	}
}
