package repository

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model/organization_schema"
)

// AccountRepository provides a interface on db level for user
type AccountRepository interface {
	CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID, tenantCode string) (int, error)
	// SelectSignInColumns select columns to generate token when user sign in
	SelectSignInColumns(uid string) (*model.SignIn, error)
	// GetOrganizationDetail get organization detail
	GetOrganizationDetail(organizationUuid string) (*model.Organization, error)
	// GetUserAccountDetail get accounts detail by uid provided
	GetUserAccountDetail(uid string) (*model.Account, error)
	// Update user language preference
	UpdateUserLanguagePreference(db *sql.DB, userId int, languagePreference string) error
	// Update user theme preference
	UpdateUserThemePreference(db *sql.DB, userId int, themePreference string) error
	// UpdateProfilePhotoOfUsers update profile_photo column in users
	UpdateProfilePhotoOfUsers(db *sql.DB, userId int, profilePhoto string) error
	// GetUserDetail get user details from organization_schema(uuid)
	GetUserDetail(db *sql.DB, userId int) (*organization_schema.User, error)
}

type accountRepositoryHandler struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepositoryHandler{db}
}
