package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

/*
Note
*/

type TransferList treegrid.GridList
type Transfer treegrid.MainRow

type (
	UploadService interface {
		Handle(req *treegrid.PostRequest, accountID int) (*treegrid.PostResponse, error)
	}
)
