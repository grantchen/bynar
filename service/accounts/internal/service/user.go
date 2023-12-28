package service

import (
	"context"
	"database/sql"
	stderr "errors"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// GipKeys db to gip key
var GipKeys = map[string]string{
	"email":     "email",
	"full_name": "displayName",
	"phone":     "phoneNumber",
	"status":    "disableUser",
}

type UserService struct {
	db                           *sql.DB
	accountDB                    *sql.DB
	organizationID               int
	customerID                   string
	authProvider                 gip.AuthProvider
	paymentProvider              checkout.PaymentClient
	simpleOrganizationRepository treegrid.SimpleGridRowRepository
	language                     string
}

func NewUserService(db *sql.DB, accountDB *sql.DB, organizationID int, customerID string, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient, simpleOrganizationService treegrid.SimpleGridRowRepository, language string) *UserService {
	return &UserService{db, accountDB, organizationID, customerID, authProvider, paymentProvider, simpleOrganizationService, language}
}

// Handle implements treegrid.Service
func (s *UserService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	resp := treegrid.HandleSingleTreegridRows(grList, func(gr treegrid.GridRow) error {
		err = utils.WithTransaction(s.db, func(tx *sql.Tx) error {
			return s.handle(tx, gr)
		})
		return i18n.TranslationErrorToI18n(s.language, err)
	})

	return resp, nil
}

// GetPageCount implements treegrid.Service
func (s *UserService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.simpleOrganizationRepository.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *UserService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {

	return s.simpleOrganizationRepository.GetPageData(tr)
}

func (s *UserService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error
	fieldsValidating := []string{"email", "phone"}

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		for _, i := range gr.UpdatedFields() {
			if i == "phone" {
				phone := strings.TrimSpace(gr["phone"].(string))
				if len(phone) > 0 && phone[0] != '+' {
					gr["phone"] = "+" + phone
				}
			}
		}
		err1 := gr.ValidateOnRequiredAll(repository.UserFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnLimitLength(repository.UserFieldNamesString, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := s.simpleOrganizationRepository.ValidateOnIntegrity(tx, gr, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
			}
		}
		err = func() error {
			err = s.simpleOrganizationRepository.Add(tx, gr)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			// create user in gip
			email, _ := gr.GetValString("email")
			fullName, _ := gr.GetValString("full_name")
			phone, _ := gr.GetValString("phone")
			status, _ := gr.GetValInt("status")
			if phone[0] != '+' {
				phone = "+" + phone
			}
			uid, err := s.authProvider.CreateUser(context.Background(), email, fullName, phone, status == 0)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			var userID int
			stmt, err := tx.Prepare("SELECT id FROM users WHERE email=?")
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			err = stmt.QueryRow(email).Scan(&userID)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			// insert into organization_accounts
			stmt, err = s.accountDB.Prepare(`INSERT INTO organization_accounts (organization_id, organization_user_uid, organization_user_id, oraginzation_main_account) VALUES(?, ?, ?, ?)`)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			_, err = stmt.Exec(s.organizationID, uid, userID, 0)
			return i18n.TranslationErrorToI18n(s.language, err)
		}()
	case treegrid.GridRowActionChanged:
		for _, i := range gr.UpdatedFields() {
			if i == "phone" {
				phone := strings.TrimSpace(gr["phone"].(string))
				if len(phone) > 0 && phone[0] != '+' {
					gr["phone"] = "+" + phone
				}
			}
		}
		err1 := gr.ValidateOnRequired(repository.UserFieldNames, s.language)
		if err1 != nil {
			return err1
		}
		err = gr.ValidateOnLimitLength(repository.UserFieldNamesString, 100, s.language)
		if err != nil {
			return err
		}
		err = func() error {
			id, ok := gr.GetValInt("id")
			if ok {
				for _, field := range fieldsValidating {
					ok, err = s.simpleOrganizationRepository.ValidateOnIntegrity(tx, gr, []string{field})
					if !ok || err != nil {
						templateData := map[string]string{
							"Field": field,
						}
						return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
					}
				}
				err = s.simpleOrganizationRepository.Update(tx, gr)
				if err != nil {
					return i18n.TranslationErrorToI18n(s.language, err)
				}

				var uid string
				stmt, err := s.accountDB.Prepare(`SELECT organization_user_uid FROM organization_accounts WHERE organization_id = ? AND organization_user_id = ?`)
				if err != nil {
					return i18n.TranslationI18n(s.language, "NoUserFound", map[string]string{})
				}
				err = stmt.QueryRow(s.organizationID, id).Scan(&uid)
				if err != nil {
					return i18n.TranslationI18n(s.language, "GipUserNotFound", map[string]string{})
				}
				// update user claims in gip
				params := map[string]interface{}{}
				customClaims := map[string]interface{}{}
				for _, i := range gr.UpdatedFields() {
					if i != "reqID" && i != "policies" {
						key, ok := GipKeys[i]
						if ok {
							if i == "status" {
								status, _ := gr.GetValInt(i)
								params[key] = status == 0
							} else {
								params[key], _ = gr.GetValString(i)
							}

						} else {
							customClaims[i], _ = gr.GetValString(i)
						}
					}
				}
				u, err := s.authProvider.GetUser(context.Background(), uid)
				if err != nil {
					return i18n.TranslationErrorToI18n(s.language, err)
				}
				if u.CustomClaims == nil {
					u.CustomClaims = map[string]interface{}{}
				}
				for k, v := range customClaims {
					u.CustomClaims[k] = v
				}
				params["customClaims"] = u.CustomClaims
				err = s.authProvider.UpdateUser(context.Background(), uid, params)
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			return nil
		}()
	case treegrid.GridRowActionDeleted:
		err = func() error {
			id, _ := gr.GetValInt("id")
			var email string
			stmt, err := tx.Prepare(`SELECT email FROM users WHERE id=?`)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			err = stmt.QueryRow(id).Scan(&email)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			// delete user in gip
			err = s.authProvider.DeleteUserByEmail(context.Background(), email)
			if err != nil {
				if stderr.Is(err, gip.ErrUserNotFound) {
					logrus.Error("delete user by email from gip ", email, err)
				} else {
					return i18n.TranslationErrorToI18n(s.language, err)
				}
			}
			_, err = s.db.Exec("DELETE FROM user_group_lines WHERE user_id = ?", id)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			if s.customerID != "" {
				var mainAccount int
				_ = s.accountDB.QueryRow("SELECT oraginzation_main_account FROM organization_accounts WHERE organization_id = ? AND organization_user_id = ?", s.organizationID, id).Scan(&mainAccount)
				if mainAccount > 0 {
					resp, err := s.paymentProvider.FetchCustomerDetails(s.customerID)
					if err == nil {
						for _, i := range resp.Instruments {
							_ = s.paymentProvider.DeleteCard(i.ID)
						}
					}
				}
			}
			err = s.simpleOrganizationRepository.Delete(tx, gr)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
			return err
		}()

	default:
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
