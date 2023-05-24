package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type UploadService struct {
	db                          *sql.DB
	updateGRUserGroupRepository treegrid.GridRowRepositoryWithChild
	updateGRUserRepository      treegrid.SimpleGridRowRepository
}

func NewUploadService(db *sql.DB,
	updateGridRowRepository treegrid.GridRowRepositoryWithChild,
	updateUserRepository treegrid.SimpleGridRowRepository,
) *UploadService {
	return &UploadService{
		db:                          db,
		updateGRUserGroupRepository: updateGridRowRepository,
		updateGRUserRepository:      updateUserRepository,
	}
}

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	b, _ := json.Marshal(req)
	logger.Debug("request: ", string(b))
	trList, err := treegrid.ParseRequestUpload(req, u.updateGRUserGroupRepository)

	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	m := make(map[string]interface{}, 0)
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
		resp.Changes = append(resp.Changes, m)

		for k := range tr.Items {
			resp.Changes = append(resp.Changes, tr.Items[k])
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Items[k]))
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

	if err := s.saveUserGroupLine(tx, tr, tr.Fields.GetID()); err != nil {
		return fmt.Errorf("save user group line: [%w]", err)
	}

	return nil
}

func (s *UploadService) saveUserGroup(tx *sql.Tx, tr *treegrid.MainRow) error {

	return s.updateGRUserGroupRepository.SaveMainRow(tx, tr)
}

func (s *UploadService) saveUserGroupLine(tx *sql.Tx, tr *treegrid.MainRow, parentID interface{}) error {
	for _, item := range tr.Items {
		logger.Debug("save group line: ", tr, "parentID: ", parentID)
		userId, err := s.getUserIdFromUserGroupLineId(tx, item.GetLineID())

		logger.Debug("userID: ", userId)
		// mark id of parent to corect field
		// item["parent"] = parentID

		// if action is update or delete, need a valid user_id
		if err != nil && (item.GetActionType() == treegrid.GridRowActionChanged ||
			item.GetActionType() == treegrid.GridRowActionDeleted) {
			return fmt.Errorf("not valid user_id: [%d] err: [%w]", userId, err)
		}

		userGroupId := item.GetID()

		// action user first then user_group_lines
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			logger.Debug("add child row")
			//add to users
			err := s.updateGRUserRepository.Add(tx, item)
			if err != nil {
				return fmt.Errorf("add child user error: [%w]", err)
			}

			// get id of user and assign to user_id
			item["user_id"] = item.GetID()

			err = s.updateGRUserGroupRepository.SaveLineAdd(tx, item)
			if err != nil {
				return fmt.Errorf("add child user groups line error: [%w]", err)
			}
		case treegrid.GridRowActionChanged:

			// assign userId to id to update user table
			item["id"] = userId

			// update user
			err := s.updateGRUserRepository.Update(tx, item)
			if err != nil {
				return fmt.Errorf("update child user error: [%w]", err)
			}

			// nothing else to update user_group_lines now, re-assign id
			item["id"] = userGroupId
		case treegrid.GridRowActionDeleted:
			logger.Debug("delete child")
			item["id"] = userId

			err := s.updateGRUserRepository.Delete(tx, item)

			// delete userid first
			if err != nil {
				return fmt.Errorf("delete child user error: [%w]", err)
			}

			// re-assign user_group_lines id
			item["id"] = userGroupId
			err = s.updateGRUserGroupRepository.SaveLineDelete(tx, item)
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
