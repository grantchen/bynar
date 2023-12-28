package service

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/repository"
)

// UploadService is the service for upload
type UploadService struct {
	db                                   *sql.DB
	updateGRUserGroupRepository          treegrid.SimpleGridRowRepository
	updateGRUserGroupRepositoryWithChild treegrid.GridRowRepositoryWithChild
	updateGRUserRepository               treegrid.SimpleGridRowRepository
	language                             string
}

// NewUploadService returns a new upload service
func NewUploadService(db *sql.DB,
	updateGRUserGroupRepository treegrid.SimpleGridRowRepository,
	updateGRUserGroupRepositoryWithChild treegrid.GridRowRepositoryWithChild,
	updateUserRepository treegrid.SimpleGridRowRepository,
	language string,
) *UploadService {
	return &UploadService{
		db:                                   db,
		updateGRUserGroupRepository:          updateGRUserGroupRepository,
		updateGRUserGroupRepositoryWithChild: updateGRUserGroupRepositoryWithChild,
		updateGRUserRepository:               updateUserRepository,
		language:                             language,
	}
}

// Handle handles the upload request
func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	trList, err := treegrid.ParseRequestUpload(req, u.updateGRUserGroupRepositoryWithChild)
	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	resp := treegrid.HandleTreegridWithChildRows(
		trList,
		func(mr *treegrid.MainRow) error {
			err = utils.WithTransaction(u.db, func(tx *sql.Tx) error {
				if err = u.saveUserGroup(tx, mr); err != nil {
					return i18n.TranslationI18n(u.language, "SaveUserGroup", map[string]string{
						"Message": err.Error(),
					})
				}
				return nil
			})
			return i18n.TranslationErrorToI18n(u.language, err)
		},
		func(mr *treegrid.MainRow, item treegrid.GridRow) error {
			err = utils.WithTransaction(u.db, func(tx *sql.Tx) error {
				if err = u.saveUserGroupLine(tx, item); err != nil {
					userGroupLineTemplateData := map[string]string{
						"Message": err.Error(),
					}
					return i18n.TranslationI18n(u.language, "SaveUserGroupLine", userGroupLineTemplateData)
				}
				return nil
			})
			return i18n.TranslationErrorToI18n(u.language, err)
		},
	)

	return resp, nil
}

// saveUserGroup saves user group
func (u *UploadService) saveUserGroup(tx *sql.Tx, tr *treegrid.MainRow) error {
	fieldsValidating := []string{"code"}

	var err error
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = tr.Fields.ValidateOnRequiredAll(repository.UserGroupFieldNames, u.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnLimitLength(repository.UserGroupFieldNames, 100, u.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := u.updateGRUserGroupRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}
		}
	case treegrid.GridRowActionChanged:
		err = tr.Fields.ValidateOnRequired(repository.UserGroupFieldNames, u.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnLimitLength(repository.UserGroupFieldNames, 100, u.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := u.updateGRUserGroupRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}
		}
	case treegrid.GridRowActionDeleted:
		// ignore id start with CR
		idStr := tr.Fields.GetIDStr()
		if !strings.HasPrefix(idStr, "CR") {
			stmt, err := tx.Prepare("DELETE FROM user_group_lines WHERE parent_id = ?")
			if err != nil {
				return err
			}

			defer func(stmt *sql.Stmt) {
				_ = stmt.Close()
			}(stmt)

			_, err = stmt.Exec(idStr)
			if err != nil {
				return i18n.TranslationErrorToI18n(u.language, err)
			}
		}
	}

	return u.updateGRUserGroupRepositoryWithChild.SaveMainRow(tx, tr)
}

// saveUserGroupLine saves user group lines
func (u *UploadService) saveUserGroupLine(tx *sql.Tx, item treegrid.GridRow) error {
	parentID := item.GetParentID()
	var err error
	switch item.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = item.ValidateOnRequiredAll(map[string][]string{"user_id": repository.UserGroupLineFieldNames["user_id"]}, u.language)
		if err != nil {
			return err
		}

		userId := item["user_id"]
		exists, err := u.checkValidUser(tx, userId)
		if err != nil {
			return fmt.Errorf("check valid user error: [%w]", err)
		}
		if !exists {
			templateData := map[string]string{
				"UserId": fmt.Sprintf("%s", userId),
			}
			return i18n.TranslationI18n(u.language, "UserNotExist", templateData)
		}

		exists, err = u.userExistInLine(tx, parentID, userId)
		if err != nil {
			return fmt.Errorf("check user exist in line error: [%w]", err)
		}
		if exists {
			templateData := map[string]string{
				"UserId": fmt.Sprintf("%s", userId),
			}
			return i18n.TranslationI18n(u.language, "UserBelongSpecificUserGroupLines", templateData)
		}

		err = u.updateGRUserGroupRepositoryWithChild.SaveLineAdd(tx, item)
		if err != nil {
			return fmt.Errorf("add child user groups line error: [%w]", err)
		}
	case treegrid.GridRowActionChanged:
		// DO NOTHING WITH ACTION UPDATE, NOT ALLOW UPDATE LINES TABLE
		return i18n.TranslationI18n(u.language, "UpdateChildLine", map[string]string{})
	case treegrid.GridRowActionDeleted:
		// re-assign user_group_lines id
		item["id"] = item.GetID()
		err = u.updateGRUserGroupRepositoryWithChild.SaveLineDelete(tx, item)
		if err != nil {
			return fmt.Errorf("delete child user group line error: [%w]", err)
		}
	default:
		return fmt.Errorf("undefined row type: %s", item.GetActionType())

	}
	return nil
}

// checkValidUser checks if user is valid
func (u *UploadService) checkValidUser(_ *sql.Tx, userId interface{}) (bool, error) {
	return utils.CheckExist(
		u.db,
		"SELECT 1 FROM users WHERE id = ?",
		userId,
	)
}

// userExistInLine checks if user exist in line
func (u *UploadService) userExistInLine(tx *sql.Tx, userGroupID, userId interface{}) (bool, error) {
	return utils.CheckExistInTx(
		tx,
		"SELECT 1 FROM user_group_lines WHERE parent_id = ? AND user_id = ?",
		userGroupID, userId,
	)
}
