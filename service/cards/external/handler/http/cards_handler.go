package http_handler

import (
	"database/sql"

	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
)

type CardHandler struct {
	cs service.CardService
}

func NewHTTPHandler(db *sql.DB, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient) *CardHandler {
	cs := service.NewCardService(db, authProvider, paymentProvider)
	return &CardHandler{cs}
}

// ListCards endpoint
func (h *CardHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		render.MethodNotAllowed(w)
		return
	}
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	resp, err := h.cs.ListCards(reqContext.Claims.AccountId)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, resp)
}

// AddCard endpoint
func (h *CardHandler) AddCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	var req model.AddCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	err = h.cs.AddCard(&models.ValidateCardRequest{ID: reqContext.Claims.AccountId, Token: req.Token, Name: req.Name, Email: req.Email})
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, nil)
}

// UpdateCard endpoint
func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	var req model.UpdateCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	h.cs.UpdateCard(reqContext.Claims.AccountId, req.SourceID)
	render.Ok(w, nil)
}

// DeleteCard endpoint
func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	var req model.DeleteCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	h.cs.DeleteCard(reqContext.Claims.AccountId, req.SourceID)
	render.Ok(w, nil)
}
