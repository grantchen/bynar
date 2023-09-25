/**
    @author: dongjs
    @date: 2023/9/13
    @description: Organization is struct of table Organizations
**/

package model

// Organization struct of table Organizations
type Organization struct {
	ID               int
	Description      string
	VatNumber        string
	Country          string
	DataSovereignty  string
	OrganizationUuid string
	TenantId         int
	Status           bool
	Verified         bool
}
