package model

type Site struct {
	Id                       int    `json:"id"`
	Code                     string `json:"code"`
	Description              string `json:"description"`
	TransactionCode          string `json:"transaction_code"`
	SubsidiariesUuid         string `json:"subsidiaries_uuid"`
	AddressUuid              string `json:"address_uuid"`
	ContactUuid              string `json:"contact_uuid"`
	ResponsibilityCenterUuid string `json:"responsibility_center_uuid"`
}

type SiteChange struct {
	Site
	Added   int `json:"Added,omitempty"`
	Changed int `json:"Changed,omitempty"`
	Deleted int `json:"Deleted,omitempty"`
}
