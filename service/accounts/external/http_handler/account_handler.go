package http_handler

import (
	"database/sql"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
)

type AccountHandler struct {
	as service.AccountService
}

func NewHTTPHandler(db *sql.DB, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient) *AccountHandler {
	as := service.NewAccountService(db, authProvider, paymentProvider)
	return &AccountHandler{as}
}

func (h *AccountHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.SignupRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	if err := h.as.Signup(req.Email); err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, nil)
}

func (h *AccountHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.ConfirmEmailRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	id, err := h.as.ConfirmEmail(req.Email, req.Timestamp, req.Signature)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, model.ConfirmEmailResponse{AccountID: id})
}

func (h *AccountHandler) VerifyCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.VerifyCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	customerID, sourceID, err := h.as.VerifyCard(req.Token, req.Email, req.Name)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, &model.VerifyCardResponse{CustomerID: customerID, SourceID: sourceID})
}

func (h *AccountHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.CreateUserRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	if !req.IsAgreementSigned {
		render.Error(w, "Agreement not signed")
		return
	}
	uid, err := h.as.CreateUser(req.Username, req.Code, req.Sign, req.Token, req.FullName, req.Country, req.AddressLine, req.AddressLine2, req.City, req.PostalCode, req.State, req.PhoneNumber, req.OrganizationName, req.VAT, req.OrganisationCountry, req.CustomerID, req.SourceID)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, model.CreateUserResponse{Token: uid})
}

// SignIn user sign in api
func (h *AccountHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.SignInRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	idToke, err := h.as.SignIn(req.Email, req.OobCode)
	if err != nil {
		render.Error(w, err.Error())
	}
	render.Ok(w, model.SignInResponse{IdToke: idToke})
}

// SendSignInEmail send sign in email of Google Identify Platform
func (h *AccountHandler) SendSignInEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.SendSignInEmailRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	err := h.as.SendSignInEmail(req.Email)
	if err != nil {
		render.Error(w, err.Error())
	}
	render.Ok(w, nil)
}
