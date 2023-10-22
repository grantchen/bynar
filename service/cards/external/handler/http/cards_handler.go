package http_handler

import (
	"database/sql"

	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
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
	if !reqContext.Claims.OrganizationAccount {
		i18n.Localize(reqContext.Claims.Language, errors.FromError(err).Code)
		render.Error(w, "not organization account")
		return
	}
	resp, err := h.cs.ListCards(reqContext.Claims.AccountId)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.FromError(err).Code))
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
	if !reqContext.Claims.OrganizationAccount {
		render.Error(w, "not organization account")
		return
	}
	var req model.AddCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	err = h.cs.AddCard(&models.ValidateCardRequest{ID: reqContext.Claims.AccountId, Token: req.Token, Name: reqContext.Claims.Name, Email: reqContext.Claims.Email})
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.FromError(err).Code))
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
	if !reqContext.Claims.OrganizationAccount {
		render.Error(w, "not organization account")
		return
	}
	var req model.UpdateCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	err = h.cs.UpdateCard(reqContext.Claims.AccountId, req.SourceID)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.FromError(err).Code))
		return
	}
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
	if !reqContext.Claims.OrganizationAccount {
		render.Error(w, "not organization account")
		return
	}
	var req model.DeleteCardRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}
	err = h.cs.DeleteCard(reqContext.Claims.AccountId, req.SourceID)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.FromError(err).Code))
		return
	}
	render.Ok(w, nil)
}
