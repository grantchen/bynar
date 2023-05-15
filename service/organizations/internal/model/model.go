package model

type Organization struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	VatNumber string `json:"vat_number"`
	State     string `json:"state"`
	Code      string `json:"code"`
}

type OrganizationChange struct {
	Organization
	Added   int `json:"Added,omitempty"`
	Changed int `json:"Changed,omitempty"`
	Deleted int `json:"Deleted,omitempty"`
}
