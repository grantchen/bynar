package service

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// UploadService is the service for upload
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

// NewUploadService returns a new upload service
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

// Handle handles the upload request
func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	trList, err := treegrid.ParseRequestUpload(req, u.grTransferRepositoryWithChild)
	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	resp := treegrid.HandleMainRowsLinesWithChild(
		trList,
		func(mr *treegrid.MainRow) error {
			err = utils.WithTransaction(u.db, func(tx *sql.Tx) error {
				return u.handle(tx, mr)
			})
			return i18n.TranslationErrorToI18n(u.language, err)
		},
	)

	return resp, nil
}

// handle handles upload request of single row
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
