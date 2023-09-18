/**
    @author: dongjs
    @date: 2023/9/7
    @description:
**/

package model

// SignIn generate token columns in database when user sign in
type SignIn struct {
	Uid                  string `json:"uid"`
	OrganizationUuid     string `json:"organization_uuid"`
	OrganizationUserId   int    `json:"organization_user_id"`
	OrganizationStatus   bool   `json:"organization_status"`
	OrganizationAccount  bool   `json:"organization_account"`
	OrganizationVerified bool   `json:"organization_verified"`
	TenantUuid           string `json:"tenant_uuid"`
	TenantStatus         bool   `json:"tenant_status"`
	TenantSuspended      bool   `json:"tenant_suspended"`

	Language string `json:"language"`
	Theme    string `json:"theme"`
}
