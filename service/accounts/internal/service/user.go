package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"log"
	"regexp"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// db to gip key
var GIP_KEYS = map[string]string{
	"email":     "email",
	"full_name": "displayName",
	"phone":     "phoneNumber",
	"status":    "disableUser",
}

type UserService struct {
	db                           *sql.DB
	accountDB                    *sql.DB
	organizationID               int
	authProvider                 gip.AuthProvider
	simpleOrganizationRepository treegrid.SimpleGridRowRepository
	language                     string
}

func NewUserService(db *sql.DB, accountDB *sql.DB, organizationID int, authProvider gip.AuthProvider, simpleOrganizationService treegrid.SimpleGridRowRepository, language string) *UserService {
	return &UserService{db, accountDB, organizationID, authProvider, simpleOrganizationService, language}
}

// Handle implements treegrid.TreeGridService
func (s *UserService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{Changes: []map[string]interface{}{}}
	// Create new transaction
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: [%w]", i18n.Localize(s.language, errors.ErrCodeBeginTransaction), err)
	}
	defer tx.Rollback()
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}
	isCommit := true
	seenEmails := make(map[string]bool)
	seenPhone := make(map[string]bool)
	for _, gr := range grList {
		if gr["email"] != nil {
			email := gr["email"].(string)
			//Check if the value is already in the map
			if seenEmails[email] {
				//If there is the same email, handle it accordingly.
				isCommit = false
				resp.IO.Result = -1
				resp.IO.Message = fmt.Sprintf("email: %s: %s", i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["email"].(string))
				resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			} else {
				seenEmails[email] = true
			}
		}
	}
	for _, gr := range grList {
		if gr["phone"] != nil {
			phone := gr["phone"].(string)
			//Check if the value is already in the map
			if seenPhone[phone] {
				//If there is the same phone, handle it accordingly.
				isCommit = false
				resp.IO.Result = -1
				resp.IO.Message = fmt.Sprintf("phone: %s: %s", i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["phone"].(string))
				resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			} else {
				seenPhone[phone] = true
			}
		}
	}
	// If no errors occurred, commit the transaction
	if isCommit == true {
		for _, gr := range grList {
			if err = s.handle(tx, gr); err != nil {
				log.Println("Err", err)
				isCommit = false
				resp.IO.Result = -1
				resp.IO.Message += err.Error() + "\n"
				resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
				break
			}
			resp.Changes = append(resp.Changes, gr)
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
		}
	}
	if isCommit == true {
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("%s: [%w]", i18n.Localize(s.language, errors.ErrCodeCommitTransaction), err)
		}
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

func (s *UserService) handle(tx *sql.Tx, gr treegrid.GridRow) error {
	var err error
	fieldsValidating := []string{"email"}
	fieldsValidatingPhone := []string{"phone"}
	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.UserFieldNames)
		if err1 != nil {
			return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeRequiredFieldsBlank))
		}
		ok, err1 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("%s: %s: %s", strings.Join(fieldsValidating, ", "), i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["email"])
		}
		ok1, err1 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidatingPhone)
		if !ok1 || err1 != nil {
			return fmt.Errorf("%s: %s: %s", strings.Join(fieldsValidatingPhone, ", "), i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["phone"])
		}
		err = func() error {
			err = s.simpleOrganizationRepository.Add(tx, gr)
			if err != nil {
				//Formatted messy string
				contains := strings.Contains(err.Error(), "too long")
				contains1 := strings.Contains(err.Error(), "Too Long")
				if contains || contains1 {
					return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeTooLong))
				} else {
					return err
				}
			}
			// create user in gip
			email, _ := gr.GetValString("email")
			fullName, _ := gr.GetValString("full_name")
			phone, _ := gr.GetValString("phone")
			status, _ := gr.GetValInt("status")
			uid, err := s.authProvider.CreateUser(context.Background(), email, fullName, phone, status == 0)
			if err != nil {
				phonePattern := `(?i)INVALID_PHONE_NUMBER|phone number`
				regex := regexp.MustCompile(phonePattern)
				emailPattern := `(?i)INVALID_EMAIL|email`
				regexEmail := regexp.MustCompile(emailPattern)
				if regex.MatchString(err.Error()) {
					return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodePhoneNumber))
				} else if regexEmail.MatchString(err.Error()) {
					return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeEmail))
				} else {
					return err
				}
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
			return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeRequiredFieldsBlank))
		}
		err = func() error {
			id, ok := gr.GetValInt("id")
			if ok {
				ok, err1 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidating)
				if !ok || err1 != nil {
					return fmt.Errorf("%s: %s: %s", strings.Join(fieldsValidating, ", "), i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["email"])
				}
				ok1, err2 := s.simpleOrganizationRepository.ValidateOnIntegrity(gr, fieldsValidatingPhone)
				if !ok1 || err2 != nil {
					return fmt.Errorf("%s: %s: %s", strings.Join(fieldsValidatingPhone, ", "), i18n.Localize(s.language, errors.ErrCodeValueDuplicated), gr["phone"])
				}
				err = s.simpleOrganizationRepository.Update(tx, gr)
				if err != nil {
					//Formatted messy string
					contains := strings.Contains(err.Error(), "too long")
					contains1 := strings.Contains(err.Error(), "Too Long")
					if contains || contains1 {
						return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeTooLong))
					} else {
						return err
					}
				}

				var uid string
				stmt, err := s.accountDB.Prepare(`SELECT organization_user_uid FROM organization_accounts WHERE organization_id = ? AND organization_user_id = ?`)
				if err != nil {
					return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeGipUser))
				}
				err = stmt.QueryRow(s.organizationID, id).Scan(&uid)
				if err != nil {
					return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeGipUser))
				}
				// update user claims in gip
				params := map[string]interface{}{}
				customClaims := map[string]interface{}{}
				for _, i := range gr.UpdatedFields() {
					if i != "reqID" {
						key, ok := GIP_KEYS[i]
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
				params["customClaims"] = customClaims
				err = s.authProvider.UpdateUser(context.Background(), uid, params)
				if err != nil {
					phonePattern := `(?i)INVALID_PHONE_NUMBER|phone number`
					regex := regexp.MustCompile(phonePattern)
					emailPattern := `(?i)INVALID_EMAIL|email`
					regexEmail := regexp.MustCompile(emailPattern)
					contains := strings.Contains(err.Error(), "user not found")
					if regex.MatchString(err.Error()) {
						return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodePhoneNumber))
					} else if regexEmail.MatchString(err.Error()) {
						return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeEmail))
					} else if contains {
						return fmt.Errorf(i18n.Localize(s.language, errors.ErrCodeGipUser))
					} else {
						return err
					}
				}
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
				contains := strings.Contains(err.Error(), "user not found")
				if !contains {
					return err
				}
			}
			err = s.simpleOrganizationRepository.Delete(tx, gr)
			return err
		}()

	default:
		return fmt.Errorf("%s: %s", i18n.Localize(s.language, errors.ErrCodeUndefinedTowType), gr.GetActionType())
	}

	if err != nil {
		return err
	}

	return err
}
