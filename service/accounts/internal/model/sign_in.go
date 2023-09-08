/**
    @author: dongjs
    @date: 2023/9/7
    @description:
**/

package model

type SignIn struct {
	Uid                     string `json:"uid"`
	OrganizationUuid        string `json:"organization_uuid"`
	OrganizationUserId      string `json:"organization_user_id"`
	OrganizationStatus      bool   `json:"organization_status"`
	TenantUuid              string `json:"tenant_uuid"`
	OrganizationMainAccount bool   `json:"organization_account"`
}
