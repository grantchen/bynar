package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"github.com/sirupsen/logrus"
)

var GIP_KEYS = map[string]string{
	"full_name": "displayName",
	"phone":     "phoneNumber",
}

type UserService struct {
	db                           *sql.DB
	accountDB                    *sql.DB
	organizationID               int
	authProvider                 gip.AuthProvider
	simpleOrganizationRepository treegrid.SimpleGridRowRepository
}

func NewUserService(db *sql.DB, accountDB *sql.DB, organizationID int, authProvider gip.AuthProvider, simpleOrganizationService treegrid.SimpleGridRowRepository) *UserService {
	return &UserService{db, accountDB, organizationID, authProvider, simpleOrganizationService}
}

func (s *UserService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}
	for _, gr := range grList {
		if err := s.handle(gr); err != nil {
			logrus.Error("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			break
		}
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}

	return resp, nil
}

func (s *UserService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.simpleOrganizationRepository.GetPageCount(tr)
	return float64(count), err
}

func (s *UserService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {

	return s.simpleOrganizationRepository.GetPageData(tr)
}

func (s *UserService) handle(gr treegrid.GridRow) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	fieldsValidating := []string{"email"}

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.UserFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%v], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = func() error {
			err = s.simpleOrganizationRepository.Add(tx, gr)
			if err != nil {
				return err
			}
			email, _ := gr.GetValString("email")
			fullName, _ := gr.GetValString("full_name")
			phone, _ := gr.GetValString("phone")
			uid, err := s.authProvider.CreateUser(context.Background(), email, fullName, phone)
			if err != nil {
				return err
			}
			var userID int
			err = tx.QueryRow(`SELECT id FROM users WHERE email=?`, email).Scan(&userID)
			if err != nil {
				return err
			}
			_, err = s.accountDB.Exec(`INSERT INTO organization_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES(?, ?, ?, ?)`, s.organizationID, uid, userID, 0)
			return err
		}()
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.UserFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%w], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = func() error {
			err = s.simpleOrganizationRepository.Update(tx, gr)
			if err != nil {
				return err
			}
			id, _ := gr.GetValInt("id")
			var email string
			err = tx.QueryRow(`SELECT email FROM users WHERE id=?`, id).Scan(&email)
			if err != nil {
				return err
			}
			params := map[string]interface{}{}
			for _, i := range gr.UpdatedFields() {
				if i != "reqID" {
					params[GIP_KEYS[i]], _ = gr.GetValString(i)
				}
			}
			err = s.authProvider.UpdateUserByEmail(context.Background(), email, params)
			return err
		}()
	case treegrid.GridRowActionDeleted:
		err = func() error {
			id, _ := gr.GetValInt("id")
			var email string
			err = tx.QueryRow(`SELECT email FROM users WHERE id=?`, id).Scan(&email)
			if err != nil {
				return err
			}
			err = s.authProvider.DeleteUserByEmail(context.Background(), email)
			if err != nil {
				return err
			}
			err = s.simpleOrganizationRepository.Delete(tx, gr)
			return err
		}()

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: [%w]", err)
	}

	return err
}
