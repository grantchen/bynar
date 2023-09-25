package models

type ServicePolicy struct {
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type Policy struct {
	Services []ServicePolicy `json:"services,omitempty"`
}
