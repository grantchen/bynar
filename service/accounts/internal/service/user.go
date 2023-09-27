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

// db to gip key
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

// Handle implements treegrid.TreeGridService
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

// GetPageCount implements treegrid.TreeGridService
func (s *UserService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.simpleOrganizationRepository.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
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
			logrus.Errorf("validate duplicate: [%v], field: %s", err1, strings.Join(fieldsValidating, ", "))
			return fmt.Errorf("validate duplicate, field: %s", strings.Join(fieldsValidating, ", "))
		}
		err = func() error {
			err = s.simpleOrganizationRepository.Add(tx, gr)
			if err != nil {
				return err
			}
			// create user in gip
			email, _ := gr.GetValString("email")
			fullName, _ := gr.GetValString("full_name")
			phone, _ := gr.GetValString("phone")
			uid, err := s.authProvider.CreateUser(context.Background(), email, fullName, phone)
			if err != nil {
				return err
			}
			var userID int
			stmt, err := tx.Prepare("SELECT id FROM users WHERE email=?")
			if err != nil {
				return err
			}
			err = stmt.QueryRow(email).Scan(&userID)
			if err != nil {
				return err
			}
			// insert into organization_accounts
			stmt, err = s.accountDB.Prepare(`INSERT INTO organization_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES(?, ?, ?, ?)`)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(s.organizationID, uid, userID, 0)
			return err
		}()
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.UserFieldNames)
		if err1 != nil {
			return err1
		}
		err = func() error {
			id, ok := gr.GetValInt("id")
			if ok {
				ok, err1 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidating)
				if !ok || err1 != nil {
					return fmt.Errorf("validate duplicate: [%w], field: %s", err1, strings.Join(fieldsValidating, ", "))
				}
				err = s.simpleOrganizationRepository.Update(tx, gr)
				if err != nil {
					return err
				}

				var email string
				stmt, err := tx.Prepare(`SELECT email FROM users WHERE id=?`)
				if err != nil {
					return err
				}
				err = stmt.QueryRow(id).Scan(&email)
				if err != nil {
					return err
				}
				// update user claims in gip
				params := map[string]interface{}{}
				customClaims := map[string]interface{}{}
				for _, i := range gr.UpdatedFields() {
					if i != "reqID" {
						key, ok := GIP_KEYS[i]
						if ok {
							params[key], _ = gr.GetValString(i)
						} else {
							customClaims[i], _ = gr.GetValString(i)
						}
					}
				}
				params["customClaims"] = customClaims
				err = s.authProvider.UpdateUserByEmail(context.Background(), email, params)
				return err
			}
			return nil
		}()
	case treegrid.GridRowActionDeleted:
		err = func() error {
			id, _ := gr.GetValInt("id")
			var email string
			stmt, err := tx.Prepare(`SELECT email FROM users WHERE id=?`)
			if err != nil {
				return err
			}
			err = stmt.QueryRow(id).Scan(&email)
			if err != nil {
				return err
			}
			// delete user in gip
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
