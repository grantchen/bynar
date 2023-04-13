package model

type Language struct {
	Id            int    `json:"id"`
	Country       string `json:"country"`
	Language      string `json:"language"`
	Two_letters   string `json:"two_letters"`
	Three_letters string `json:"three_letters"`
	Number        int64  `json:"number"`
}

type Changes struct {
	Id           string `json:"id,omitempty"`
	Added        int    `json:"Added,omitempty"`
	Changed      int    `json:"Changed,omitempty"`
	Deleted      int    `json:"Deleted,omitempty"`
	Language     string `json:"language,omitempty"`
	Country      string `json:"country,omitempty"`
	TwoLetters   string `json:"two_letters,omitempty"`
	ThreeLetters string `json:"three_letters,omitempty"`
	Number       string `json:"number,omitempty"`
}

// ChangedRow: used to return Messages for POST update
type ChangedRow struct {
	Id      string `json:"id,omitempty"`
	Changed int    `json:"Changed,omitempty"`
	Added   int    `json:"Added,omitempty"`
	Deleted int    `json:"Deleted,omitempty"`
	Color   string `json:"Color,omitempty"`
}

// PostRequest struct for mapping to post requests
type PostRequest struct {
	IO struct {
		Message string
	}
	Changes []Changes
}

// Response struct for json responses
type Response struct {
	IO struct {
		Message string
		Result  int
	}
	Changes []ChangedRow
}
