/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
)

// GetOrganizationDetail query organization from database by organizationUuid
func (r *accountRepositoryHandler) GetOrganizationDetail(organizationUuid string) (*model.Organization, error) {
	var querySql = `
		select a.id,coalesce(a.description,''),coalesce(a.vat_number,''),coalesce(a.country,''),
		       coalesce(a.data_sovereignty,''),coalesce(a.organization_uuid,''),a.tenant_id,a.status,a.verified
		from organizations a where a.organization_uuid = ? limit 1`
	var organization = model.Organization{}
	prepare, err := r.db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer prepare.Close()
	err = prepare.QueryRow(organizationUuid).Scan(&organization.ID,
		&organization.Description, &organization.VatNumber,
		&organization.Country, &organization.DataSovereignty,
		&organization.OrganizationUuid, &organization.TenantId, &organization.Status, &organization.Verified)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &organization, nil
}

// GetUserAccountDetail get accounts detail by uid provided
func (r *accountRepositoryHandler) GetUserAccountDetail(email string) (*model.Account, error) {
	var querySql = `
		select a.id,a.email,a.full_name,a.address,coalesce(a.address_2,''),a.phone,a.city,a.postal_code,a.country,a.state,a.status,a.uid,a.org_id,a.verified
		from accounts a where email = ? and status = ? and verified = ? limit 1`
	var account = model.Account{}
	prepare, err := r.db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer prepare.Close()
	err = prepare.QueryRow(email, true, true).Scan(&account.ID,
		&account.Email, &account.FullName, &account.Address, &account.Address2, &account.Phone,
		&account.City, &account.PostalCode, &account.Country, &account.State,
		&account.Status, &account.UID, &account.OrgID, &account.Verified)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &account, nil
}

// Update user language preference
func (r *accountRepositoryHandler) UpdateUserLanguagePreference(db *sql.DB, email, languagePreference string) error {
	var updateSql = `update users set language_preference = ? where email = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer prepare.Close()
	if _, err = prepare.Exec(languagePreference, email); err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}

	return nil
}

// Update user theme preference
func (r *accountRepositoryHandler) UpdateUserThemePreference(db *sql.DB, email, themePreference string) error {
	var updateSql = `update users set theme = ? where email = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer prepare.Close()
	if _, err = prepare.Exec(themePreference, email); err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}

	return nil
}

// UpdateProfilePhotoOfUsers update column profile_photo in table users of organization_schema(uuid)
func (r *accountRepositoryHandler) UpdateProfilePhotoOfUsers(db *sql.DB, email string, profilePhoto string) error {
	updateSql := `UPDATE users SET profile_photo = ? WHERE email = ?`
	prepare, err := db.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("db prepare: [%w], sql string: [%s]", err, updateSql)
	}
	defer prepare.Close()
	if _, err = prepare.Exec(profilePhoto, email); err != nil {
		return fmt.Errorf("db update exec: [%w]", err)
	}
	return nil
}

// GetUserDetail get user details from organization_schema(uuid)
func (r *accountRepositoryHandler) GetUserDetail(db *sql.DB, email string) (*model.User, error) {
	var querySql = `select a.id,
       a.email,
       coalesce(a.full_name,''),
       coalesce(a.phone,''),
       a.status,
       coalesce(a.language_preference,''),
       coalesce(a.policy_id,0),
       coalesce(a.theme,''),
       coalesce(a.profile_photo,'')
		from users a
		where a.email = ? and status = ? limit 1`
	var user = model.User{}
	prepare, err := db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	defer prepare.Close()
	err = prepare.QueryRow(email, true).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Phone, &user.Status,
		&user.LanguagePreference, &user.PolicyId, &user.Theme, &user.ProfilePhoto)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &user, nil
}

// GetUserPolicy get user policy from organization_schema(uuid)
func (r *accountRepositoryHandler) GetUserPolicy(db *sql.DB, id int) (map[string]int, error) {
	params := []string{"user_list", "sales"}
	var querySql = fmt.Sprintf("select %s FROM policies WHERE id=?", strings.Join(params, ","))
	prepare, err := db.Prepare(querySql)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, querySql)
	}
	cols := make([]int, len(params))
	colsp := make([]interface{}, len(params))
	for i, _ := range cols {
		colsp[i] = &cols[i]
	}
	defer prepare.Close()
	err = prepare.QueryRow(id).Scan(colsp...)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	m := make(map[string]int)
	for i, name := range params {
		val := colsp[i].(*int)
		m[name] = *val
	}

	return m, nil
}
