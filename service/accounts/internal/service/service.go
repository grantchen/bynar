package service

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"mime/multipart"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// AccountService is a interface which provide helper methods to access account related operations
type AccountService interface {
	Signup(email string) error
	ConfirmEmail(email, timestamp, signature string) (int, error)
	VerifyCard(token, email, name string) (string, string, error)
	CreateUser(email, code, sign, token, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID string) (string, error)
	// SignIn user sign in with Google identify platform oobCode
	SignIn(email, oobCode string) (string, error)
	// SendSignInEmail send sign in email of Google identify platform
	SendSignInEmail(email string) error
	VerifyEmail(email string) error
	GetUserDetails(db *sql.DB, email string) (*model.GetUserResponse, error)
	// UploadFileToGCS upload user's profile picture to google cloud storage
	UploadFileToGCS(db *sql.DB, organizationUuid, email string, multipartReader *multipart.Reader) (string, error)
	// DeleteFileFromGCS delete user's profile picture from google cloud storage
	DeleteFileFromGCS(db *sql.DB, organizationUuid, email string) error
	// Update user language preference
	UpdateUserLanguagePreference(db *sql.DB, email, languagePreference string) error
	// Update user theme preference
	UpdateUserThemePreference(db *sql.DB, email, themePreference string) error
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
