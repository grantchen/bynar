package model

type AddCardRequest struct {
	Token string `json:"token" valid:""`
}

type UpdateCardRequest struct {
	SourceID string `json:"source_id" valid:""`
}

type DeleteCardRequest struct {
	SourceID string `json:"source_id" valid:""`
}
