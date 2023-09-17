/**
    @author: dongjs
    @date: 2023/9/7
    @description:
**/

package model

// SignIn generate token columns in database when user sign in
type SignIn struct {
	Uid                     string `json:"uid"`
	OrganizationUuid        string `json:"organization_uuid"`
	OrganizationUserId      string `json:"organization_user_id"`
	OrganizationStatus      bool   `json:"organization_status"`
	TenantUuid              string `json:"tenant_uuid"`
	OrganizationMainAccount bool   `json:"organization_account"`
	LanguagePreference      string `json:"language_preference"`
}
