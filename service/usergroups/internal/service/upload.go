package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService struct {
	db                      *sql.DB
	updateGridRowRepository treegrid.GridRowRepositoryWithChild
}

func NewUploadService(db *sql.DB,
	updateGridRowRepository treegrid.GridRowRepositoryWithChild,
) *UploadService {
	return &UploadService{
		db:                      db,
		updateGridRowRepository: updateGridRowRepository,
	}
}

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	b, _ := json.Marshal(req.Changes)
	logger.Debug("request: ", string(b))

	trList, err := treegrid.ParseRequestUpload(req, u.updateGridRowRepository)

	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	for _, tr := range trList.MainRows() {
		if err := u.handle(tr); err != nil {
			log.Println("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(tr.Fields))
			break
		}
		resp.Changes = append(resp.Changes, tr.Fields)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Fields))

		for k := range tr.Items {
			resp.Changes = append(resp.Changes, tr.Items[k])
		}
	}

	return resp, nil
}

func (s *UploadService) handle(tr *treegrid.MainRow) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	if err := s.save(tx, tr); err != nil {
		return fmt.Errorf("usergroups svc save '%s': [%w]", tr.IDString(), err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: [%w]", err)
	}

	return nil
}

func (s *UploadService) save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := s.saveUserGroup(tx, tr); err != nil {
		return fmt.Errorf("save usergroups: [%w]", err)
	}

	// if err := s.SaveTransferLines(tx, tr); err != nil {
	// 	return fmt.Errorf("save transfer line: [%w]", err)
	// }

	return nil
}

func (s *UploadService) saveUserGroup(tx *sql.Tx, tr *treegrid.MainRow) error {
	return s.updateGridRowRepository.SaveMainRow(tx, tr)
}
