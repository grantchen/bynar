package http_handler

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	i18n "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"net/http"
	"strconv"
	"strings"

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
		handler.LogInternalError(err)
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
		handler.LogInternalError(errors.NewUnknownError("decode request fail", errors.ErrCodeNoUserFound).WithInternalCause(err))
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
		handler.LogInternalError(errors.NewUnknownError("verify token fail", "").WithInternalCause(err))
		render.Error(w, err.Error())
		return
	}
	userResponse, getErr := h.as.GetUserDetails(reqContext.DynamicDB, reqContext.Claims.Uid, reqContext.Claims.OrganizationUserId)
	if getErr != nil {
		handler.LogInternalError(getErr)
		render.Error(w, getErr.Error())
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
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid))
		return
	}
	reader, err := r.MultipartReader()
	if err != nil || reader == nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCode))
		return
	}
	url, uploadErr := h.as.UploadFileToGCS(reqContext.DynamicDB, reqContext.Claims.OrganizationUuid, reqContext.Claims.OrganizationUserId, reader)
	if uploadErr != nil {
		handler.LogInternalError(uploadErr)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, uploadErr.Code))
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
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid))
		return
	}

	deleteErr := h.as.DeleteFileFromGCS(reqContext.DynamicDB, reqContext.Claims.OrganizationUuid, reqContext.Claims.OrganizationUserId)

	if deleteErr != nil {
		handler.LogInternalError(deleteErr)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, deleteErr.Code))
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
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid))
		return
	}
	updateErr := h.as.UpdateUserLanguagePreference(reqContext.DynamicDB, reqContext.Claims.Uid, reqContext.Claims.OrganizationUserId, req.LanguagePreference)
	if updateErr != nil {
		handler.LogInternalError(updateErr)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, updateErr.Code))
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
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid))
		return
	}
	updateErr := h.as.UpdateUserThemePreference(reqContext.DynamicDB, reqContext.Claims.OrganizationUserId, req.ThemePreference)
	if updateErr != nil {
		handler.LogInternalError(updateErr)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, updateErr.Code))
		return
	}

	render.Ok(w, nil)
}

// UpdateUserProfile update user profile information email phone_number language theme full_name
func (h *AccountHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize("", errors.ErrCodeIDTokenInvalid))
		return
	}

	if r.Method != http.MethodPut {
		render.ErrorWithHttpCode(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid), http.StatusMethodNotAllowed)
		return
	}

	var req model.UpdateUserProfileRequest
	if err = render.DecodeJSON(r.Body, &req); err != nil {
		render.ErrorWithHttpCode(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateErr := h.as.UpdateUserProfile(reqContext.DynamicDB, reqContext.Claims.OrganizationUserId, reqContext.Claims.Uid, req)
	if updateErr != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, updateErr.Code))
		return
	}
	render.Ok(w, nil)
}

// GetUserProfileById get user profile
func (h *AccountHandler) GetUserProfileById(w http.ResponseWriter, r *http.Request) {
	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize("", errors.ErrCodeIDTokenInvalid))
		return
	}
	if r.Method != http.MethodGet {
		render.ErrorWithHttpCode(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/user/")
	if "" == id {
		render.ErrorWithHttpCode(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeRequestParameter), http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		render.ErrorWithHttpCode(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeRequestParameter), http.StatusBadRequest)
		return
	}
	if userId != reqContext.Claims.OrganizationUserId {
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodePermissionDenied))
		return
	}
	profile, reErr := h.as.GetUserProfileById(reqContext.DynamicDB, userId)
	if reErr != nil {
		handler.LogInternalError(reErr)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, reErr.Code))
		return
	}
	render.Ok(w, profile)
}

// GetOrganizationAccount get organization account information
func (h *AccountHandler) GetOrganizationAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		render.MethodNotAllowed(w)
		return
	}

	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize("", errors.ErrCodeIDTokenInvalid))
		return
	}

	if !reqContext.Claims.OrganizationAccount {
		render.Error(w, i18n.TranslationI18n(reqContext.Claims.Language, "permission-denied", nil).Error())
		return
	}

	response, err := h.as.GetOrganizationAccount(reqContext.Claims.Language, reqContext.Claims.AccountId, reqContext.Claims.OrganizationUuid)
	if err != nil {
		render.Error(w, err.Error())
		return
	}

	render.Ok(w, response)
}

// UpdateOrganizationAccount update organization account(tenant)
func (h *AccountHandler) UpdateOrganizationAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render.MethodNotAllowed(w)
		return
	}

	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid))
		return
	}

	if !reqContext.Claims.OrganizationAccount {
		render.Error(w, i18n.TranslationI18n(reqContext.Claims.Language, "permission-denied", nil).Error())
		return
	}

	var req model.OrganizationAccountRequest
	if err = render.DecodeJSON(r.Body, &req); err != nil {
		render.ErrorWithHttpCode(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.as.UpdateOrganizationAccount(
		reqContext.DynamicDB,
		reqContext.Claims.Language,
		reqContext.Claims.AccountId,
		reqContext.Claims.Uid,
		reqContext.Claims.OrganizationUserId,
		reqContext.Claims.OrganizationUuid,
		req,
	)
	if err != nil {
		render.Error(w, err.Error())
		return
	}

	render.Ok(w, nil)
}

// DeleteOrganizationAccount delete organization account(tenant)
func (h *AccountHandler) DeleteOrganizationAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		render.MethodNotAllowed(w)
		return
	}

	reqContext, err := middleware.GetIdTokenClaimsFromHttpRequestContext(r)
	if err != nil {
		handler.LogInternalError(err)
		render.Error(w, i18n.Localize(reqContext.Claims.Language, errors.ErrCodeIDTokenInvalid))
		return
	}

	if !reqContext.Claims.OrganizationAccount {
		render.Error(w, i18n.TranslationI18n(reqContext.Claims.Language, "permission-denied", nil).Error())
		return
	}

	err = h.as.DeleteOrganizationAccount(
		reqContext.DynamicDB,
		reqContext.Claims.Language,
		reqContext.Claims.TenantUuid,
		reqContext.Claims.OrganizationUuid,
	)
	if err != nil {
		render.Error(w, err.Error())
		return
	}

	render.Ok(w, nil)
}
