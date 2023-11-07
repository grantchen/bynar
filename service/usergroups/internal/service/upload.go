package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"log"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/repository"
)

type UploadService struct {
	db                                   *sql.DB
	updateGRUserGroupRepository          treegrid.SimpleGridRowRepository
	updateGRUserGroupRepositoryWithChild treegrid.GridRowRepositoryWithChild
	updateGRUserRepository               treegrid.SimpleGridRowRepository
	language                             string
}

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

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	b, _ := json.Marshal(req)
	logger.Debug("request: ", string(b))
	trList, err := treegrid.ParseRequestUpload(req, u.updateGRUserGroupRepositoryWithChild)

	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeBeginTransaction))
	}
	defer tx.Rollback()

	m := make(map[string]interface{}, 0)
	var handleErr error
	for _, tr := range trList.MainRows() {
		if handleErr = u.handle(tx, tr); handleErr != nil {
			log.Println("Err", handleErr)

			resp.IO.Result = -1
			resp.IO.Message += handleErr.Error() + "\n"
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

	if handleErr == nil {
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, errors.ErrCodeCommitTransaction), err)
		}
	}

	return resp, nil
}

func (s *UploadService) handle(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := s.save(tx, tr); err != nil {
		return i18n.TranslationErrorToI18n(s.language, err)
	}
	return nil
}

func (s *UploadService) save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := s.saveUserGroup(tx, tr); err != nil {
		userTemplateData := map[string]string{
			"Message": err.Error(),
		}
		return i18n.TranslationI18n(s.language, "SaveUserGroup", userTemplateData)
	}

	if err := s.saveUserGroupLine(tx, tr, tr.Fields.GetID()); err != nil {
		userGroupLineTemplateData := map[string]string{
			"Message": err.Error(),
		}
		return i18n.TranslationI18n(s.language, "SaveUserGroupLine", userGroupLineTemplateData)
	}

	return nil
}

func (s *UploadService) saveUserGroup(tx *sql.Tx, tr *treegrid.MainRow) error {
	fieldsValidating := []string{"code"}

	var err error
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = tr.Fields.ValidateOnRequiredAll(repository.UserGroupFieldNames, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnLimitLength(repository.UserGroupFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := s.updateGRUserGroupRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
			}
		}
	case treegrid.GridRowActionChanged:
		err = tr.Fields.ValidateOnRequired(repository.UserGroupFieldNames, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnLimitLength(repository.UserGroupFieldNames, 100, s.language)
		if err != nil {
			return err
		}
		for _, field := range fieldsValidating {
			ok, err := s.updateGRUserGroupRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(s.language, "ValueDuplicated", templateData)
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

			defer stmt.Close()

			_, err = stmt.Exec(idStr)
			if err != nil {
				return i18n.TranslationErrorToI18n(s.language, err)
			}
		}
	}

	return s.updateGRUserGroupRepositoryWithChild.SaveMainRow(tx, tr)
}

func (s *UploadService) saveUserGroupLine(tx *sql.Tx, tr *treegrid.MainRow, parentID interface{}) error {
	for _, item := range tr.Items {
		logger.Debug("save group line: ", tr, "parentID: ", parentID)

		var err error
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err = item.ValidateOnRequiredAll(map[string][]string{"user_id": repository.UserGroupLineFieldNames["user_id"]}, s.language)
			if err != nil {
				return err
			}

			logger.Debug("add child row")
			userId := item["user_id"]
			ok, err := s.checkValidUser(tx, userId)

			if err != nil || !ok {
				templateData := map[string]string{
					"UserId": fmt.Sprintf("%s", userId),
				}
				return i18n.TranslationI18n(s.language, "UserNotExist", templateData)
			}

			ok, err = s.userExistInLine(tx, userId)

			if err != nil || !ok {
				templateData := map[string]string{
					"UserId": fmt.Sprintf("%s", userId),
				}
				return i18n.TranslationI18n(s.language, "UserBelongSpecificUserGroupLines", templateData)
			}

			err = s.updateGRUserGroupRepositoryWithChild.SaveLineAdd(tx, item)
			if err != nil {
				return fmt.Errorf("add child user groups line error: [%w]", err)
			}
		case treegrid.GridRowActionChanged:
			// DO NOTHING WITH ACTION UPDATE, NOT ALLOW UPDATE LINES TABLE
			return i18n.TranslationI18n(s.language, "UpdateChildLine", map[string]string{})
		case treegrid.GridRowActionDeleted:
			logger.Debug("delete child")

			// re-assign user_group_lines id
			item["id"] = item.GetID()
			err = s.updateGRUserGroupRepositoryWithChild.SaveLineDelete(tx, item)
			if err != nil {
				return fmt.Errorf("delete child user group line error: [%w]", err)
			}
		default:
			return fmt.Errorf("undefined row type: %s", tr.Fields.GetActionType())

		}
	}
	return nil
}

func (s *UploadService) getUserIdFromUserGroupLineId(tx *sql.Tx, userGroupLineId string) (int, error) {
	query := `SELECT user_id FROM user_group_lines WHERE id = ?`
	args := []interface{}{userGroupLineId}
	rows, err := tx.Query(query, args...)
	if err != nil {
		return 0, fmt.Errorf("query: [%w], sql string: [%s]", err, query)
	}
	defer rows.Close()
	rowVals, err := utils.NewRowVals(rows)
	if err != nil {
		return 0, fmt.Errorf("new row vals: [%w], row vals: [%v]", err, rowVals)
	}

	rows.Next()
	if err := rowVals.Parse(rows); err != nil {
		return 0, fmt.Errorf("parse rows: [%w]", err)
	}

	entry := rowVals.StringValues()
	userId, _ := strconv.Atoi(entry["user_id"])
	if err != nil {
		return 0, fmt.Errorf("parse id error: [%w]", err)
	}
	return userId, nil
}

func (s *UploadService) checkValidUser(tx *sql.Tx, userId interface{}) (bool, error) {
	query := `
	SELECT COUNT(*) as Count FROM users where id = ?
	`
	params := []interface{}{userId}
	rows, err := s.db.Query(query, params...)

	if err != nil {
		return false, err
	}
	defer rows.Close()
	count, err := utils.CheckCoutWithError(rows)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (s *UploadService) userExistInLine(tx *sql.Tx, userId interface{}) (bool, error) {
	query := `
	SELECT COUNT(*) as Count FROM user_group_lines where user_id = ?
	`
	params := []interface{}{userId}
	rows, err := s.db.Query(query, params...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	count, err := utils.CheckCoutWithError(rows)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
