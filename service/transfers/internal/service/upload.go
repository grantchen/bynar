package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService struct {
	db                            *sql.DB
	grTransferRepositoryWithChild treegrid.GridRowRepositoryWithChild
	userRepository                repository.UserRepository
	transferRepository            repository.TransferRepository
	inventoryRepository           repository.InventoryRepository
	documentRepository            repository.DocumentRepository
	accountID                     int
	language                      string
}

func NewUploadService(db *sql.DB,
	grTransferRepositoryWithChild treegrid.GridRowRepositoryWithChild,
	userRepository repository.UserRepository,
	transferRepository repository.TransferRepository,
	inventoryRepository repository.InventoryRepository,
	documentRepository repository.DocumentRepository,
	accountID int,
	language string,
) *UploadService {
	return &UploadService{
		db:                            db,
		grTransferRepositoryWithChild: grTransferRepositoryWithChild,
		userRepository:                userRepository,
		transferRepository:            transferRepository,
		inventoryRepository:           inventoryRepository,
		documentRepository:            documentRepository,
		accountID:                     accountID,
		language:                      language,
	}
}

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	trList, err := treegrid.ParseRequestUpload(req, u.grTransferRepositoryWithChild)

	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeBeginTransaction))
	}
	defer tx.Rollback()

	m := make(map[string]interface{}, 0)
	for _, tr := range trList.MainRows() {
		if err := u.handle(tx, tr); err != nil {
			log.Println("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(tr.Fields))
			break
		}
		resp.Changes = append(resp.Changes, tr.Fields)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Fields))
		resp.Changes = append(resp.Changes, m)

		for k := range tr.Items {
			resp.Changes = append(resp.Changes, tr.Items[k])
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Items[k]))
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, errors.ErrCodeCommitTransaction), err)
	}

	return resp, nil
}

func (s *UploadService) handle(tx *sql.Tx, tr *treegrid.MainRow) error {
	switch tr.Status() {
	// update/add
	case 0:
		if err := s.transferRepository.Save(tx, tr); err != nil {
			return err
		}
	case 1:
		ok, err := s.inventoryRepository.CheckQuantityAndValue(tx, tr)
		if err != nil {
			return fmt.Errorf("check inventory quantity and value: [%w]", err)
		}

		if !ok {
			return ErrInvalidQuantity
		}

		if err = s.transferRepository.Save(tx, tr); err != nil {
			return err
		}

		if err := s.inventoryRepository.Save(tx, tr); err != nil {
			return fmt.Errorf("transfer svc save: [%w]", err)
		}

		ok, err = s.documentRepository.IsAuto(tx, tr)
		if err != nil {
			return fmt.Errorf("document svc check if is auto: [%w]", err)
		}

		if !ok {
			return nil
		}

		docIdStr, err := s.documentRepository.Generate(tx, tr)
		if err != nil {
			return fmt.Errorf("document svc generate: [%w]", err)
		}

		if err := s.transferRepository.SaveDocumentID(tx, tr, docIdStr); err != nil {
			return fmt.Errorf("transfer svc save document id: [%w]", err)
		}
	default:
		if err := s.transferRepository.UpdateStatus(tx, tr.Status()); err != nil {
			return fmt.Errorf("transfer svc update status: [%w]", err)
		}
	}

	return nil
}
