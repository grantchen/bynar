package model

type Language struct {
	Id           int    `json:"id"`
	Country      string `json:"country"`
	Language     string `json:"language"`
	TwoLetters   string `json:"two_letters"`
	ThreeLetters string `json:"three_letters"`
	Number       int    `json:"number"`
}

type LanguageChange struct {
	Language
	Added   int `json:"Added,omitempty"`
	Changed int `json:"Changed,omitempty"`
	Deleted int `json:"Deleted,omitempty"`
}
