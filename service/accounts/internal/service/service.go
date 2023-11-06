package service

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"mime/multipart"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// AccountService is a interface which provide helper methods to access account related operations
type AccountService interface {
	Signup(email string) *errors.Error
	ConfirmEmail(email, timestamp, signature string) (int, *errors.Error)
	VerifyCard(token, email, name string) (string, string, *errors.Error)
	CreateUser(email, code, sign, token, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID, tenantCode string) (string, *errors.Error)
	// SignIn user sign in with Google identify platform oobCode
	SignIn(email, oobCode string) (string, *errors.Error)
	// SendSignInEmail send sign in email of Google identify platform
	SendSignInEmail(email string) *errors.Error
	VerifyEmail(email string) *errors.Error
	GetUserDetails(db *sql.DB, uid string, userId int) (*model.GetUserResponse, *errors.Error)
	// UploadFileToGCS upload user's profile picture to google cloud storage
	UploadFileToGCS(db *sql.DB, organizationUuid string, userId int, multipartReader *multipart.Reader) (string, *errors.Error)
	// DeleteFileFromGCS delete user's profile picture from google cloud storage
	DeleteFileFromGCS(db *sql.DB, organizationUuid string, userId int) *errors.Error
	// Update user language preference
	UpdateUserLanguagePreference(db *sql.DB, uid string, userId int, languagePreference string) *errors.Error
	// Update user theme preference
	UpdateUserThemePreference(db *sql.DB, userId int, themePreference string) *errors.Error
	// UpdateUserProfile update user profile
	UpdateUserProfile(db *sql.DB, userId int, uid string, userProfile model.UpdateUserProfileRequest) *errors.Error
	// UpdateGipCustomClaims update custom claims of Google Identify Platform
	UpdateGipCustomClaims(uid string) *errors.Error
	// GetUserProfileById get user profile info
	GetUserProfileById(db *sql.DB, userId int) (*model.UserProfileResponse, *errors.Error)
	// GetOrganizationAccount get organization account information
	GetOrganizationAccount(language string, accountID int, organizationUuid string) (*model.GetOrganizationAccountResponse, error)
	// UpdateOrganizationAccount update organization account
	UpdateOrganizationAccount(db *sql.DB, language string, accountID int, uid string, organizationUserId int, organizationUuid string, organizationAccount model.OrganizationAccountRequest) error
	// DeleteOrganizationAccount delete organization account
	DeleteOrganizationAccount(db *sql.DB, language string, tenantUuid string, organizationUuid string) error
}

type accountServiceHandler struct {
	ar                   repository.AccountRepository
	authProvider         gip.AuthProvider
	paymentProvider      checkout.PaymentClient
	cloudStorageProvider gcs.CloudStorageProvider
}

// NewAccountService initiates the account service object
func NewAccountService(db *sql.DB, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient, cloudStorageProvider gcs.CloudStorageProvider) AccountService {
	ar := repository.NewAccountRepository(db)
	return &accountServiceHandler{ar, authProvider, paymentProvider, cloudStorageProvider}
}
