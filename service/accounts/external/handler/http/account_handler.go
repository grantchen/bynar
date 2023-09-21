package http_handler

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	i18n "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
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

func NewHTTPHandler(db *sql.DB, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient, cloudStorageProvider gcs.CloudStorageProvider) *AccountHandler {
	as := service.NewAccountService(db, authProvider, paymentProvider, cloudStorageProvider)
	return &AccountHandler{as}
}

// Signup endpoint use gip to send email
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

// ConfirmEmail endpoint
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

// VerifyCard endpoint
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

// CreateUser endpoint
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

	phoneNumber := req.PhoneNumber
	if phoneNumber[0] != '+' {
		phoneNumber = "+" + phoneNumber
	}
	customToken, err := h.as.CreateUser(req.Username, req.Timestamp, req.Signature, req.Token, req.FullName, req.Country, req.AddressLine, req.AddressLine2, req.City, req.PostalCode, req.State, phoneNumber, req.OrganizationName, req.VAT, req.OrganisationCountry, req.CustomerID, req.SourceID, req.TenantCode)
	if err != nil {
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, model.CreateUserResponse{Token: customToken})
}

// SignIn user sign in api
func (h *AccountHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		render.MethodNotAllowed(w)
		return
	}
	var req model.SignInRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		handler.LogInternalError(err)
		render.Error(w, err.Error())
		return
	}
	idToke, err := h.as.SignIn(req.Email, req.OobCode)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, err.Error())
		return
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
		handler.LogInternalError(errors.NewUnknownError("decode request fail").WithInternalCause(err))
		render.Error(w, err.Error())
		return
	}
	err := h.as.SendSignInEmail(req.Email)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, nil)
}

// User get user info when user sign_in
func (h *AccountHandler) User(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		render.MethodNotAllowed(w)
		return
	}
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(errors.NewUnknownError("verify token fail").WithInternalCause(err))
		render.Error(w, err.Error())
		return
	}
	userResponse, err := h.as.GetUserDetails(reqContext.DynamicDB, reqContext.Claims.Email)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, err.Error())
		return
	}
	render.Ok(w, userResponse)
}

// UploadProfilePhoto upload profile_photo
func (h *AccountHandler) UploadProfilePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render.MethodNotAllowed(w)
		return
	}

	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "upload-profile-fail"))
		return
	}
	reader, err := r.MultipartReader()
	if err != nil || reader == nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "upload-profile-fail"))
		return
	}
	url, err := h.as.UploadFileToGCS(reqContext.DynamicDB, reqContext.Claims.OrganizationUuid, reqContext.Claims.Email, reader)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "upload-profile-fail"))
		return

	}
	render.Ok(w, url)
}

// DeleteProfileImage delete user's profile_picture from google cloud storage
func (h *AccountHandler) DeleteProfileImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		render.MethodNotAllowed(w)
		return
	}
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "delete-profile-fail"))
		return
	}

	err = h.as.DeleteFileFromGCS(reqContext.DynamicDB, reqContext.Claims.OrganizationUuid, reqContext.Claims.Email)

	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "delete-profile-fail"))
		return
	}
	render.Ok(w, nil)
}

// Update user language preference
func (h *AccountHandler) UpdateUserLanguagePreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		render.MethodNotAllowed(w)
		return
	}
	var req model.UpdateUserLanguagePreferenceRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}

	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "error"))
		return
	}
	err = h.as.UpdateUserLanguagePreference(reqContext.DynamicDB, reqContext.Claims.Email, req.LanguagePreference)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "error"))
		return
	}

	render.Ok(w, nil)
}

// Update user theme preference
func (h *AccountHandler) UpdateUserThemePreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		render.MethodNotAllowed(w)
		return
	}
	var req model.UpdateUserThemePreferenceRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Error(w, err.Error())
		return
	}

	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "error"))
		return
	}
	err = h.as.UpdateUserThemePreference(reqContext.DynamicDB, reqContext.Claims.Email, req.ThemePreference)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, "error"))
		return
	}

	render.Ok(w, nil)
}
