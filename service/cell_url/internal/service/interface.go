package service

import (
	"context"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/model"
)

type DataGridService interface {
	GetGridDataByValue(ctx context.Context, value string, id string) (*model.Grid, error)
}
