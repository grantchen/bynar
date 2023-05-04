package repository

import (
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type workflowRepository struct {
	conn     *sql.DB
	moduleID int
}

func NewWorkflowRepository(conn *sql.DB, moduleID int) WorkflowRepository {
	return &workflowRepository{
		conn:     conn,
		moduleID: moduleID,
	}
}

// TODO: approval order. Using parent_it but need module id

func (s *workflowRepository) GetWorkflowItem(moduleID, accountID, documentID int) (models.WorkflowItem, error) {
	query := `
	SELECT wi.id, wi.approval_order, d.status
	FROM workflow_items wi
		INNER JOIN workflows w ON w.id = wi.parent_id
		INNER JOIN documents d ON wi.document_id = d.id
	WHERE wi.document_id = ? AND wi.account_id = ? AND w.module_id = ? 
	`

	wrkItem := models.WorkflowItem{
		ParentID:   moduleID,
		DocumentID: documentID,
		AccountID:  accountID,
	}

	err := s.conn.QueryRow(query, documentID, accountID, moduleID).Scan(&wrkItem.Id, &wrkItem.ApprovalOrder, &wrkItem.Status)
	if err != nil {
		return wrkItem, fmt.Errorf("query row: [%w], query: %s, %d, %d, %d", err, query, documentID, accountID, moduleID)
	}

	return wrkItem, nil
}

func (s *workflowRepository) GetModuleID() int {
	return s.moduleID
}

// CheckApprovalOrder
// Get last transfer status, calc next one, and check if user can access to next status or
// if transfer is new:
// linhhd: not used
//
//	-
func (s *workflowRepository) CheckApprovalOrder(conn *sql.DB, tr *treegrid.MainRow, accountID int) (bool, error) {
	// with status = 0 skip check
	if tr.Status() == 0 {
		return true, nil
	}

	if tr.Fields.GetActionType() == treegrid.GridRowActionAdd {
		return s.checkActionAdded(conn, tr, accountID)
	}

	if tr.Fields.GetActionType() == treegrid.GridRowActionChanged {
		return s.checkActionUpdated(conn, tr, accountID)
	}

	return true, nil
}

func (s *workflowRepository) checkActionAdded(conn *sql.DB, tr *treegrid.MainRow, accountID int) (bool, error) {
	docID, ok := tr.Fields.GetValInt("document_id")
	if !ok {
		return false, fmt.Errorf("missing required param document_id: %v", tr.Fields["document_id"])
	}

	storeID, ok := tr.Fields.GetValInt("store_id")
	if !ok {
		return false, fmt.Errorf("missing required param store_id")
	}

	query := `
	SELECT wi.approval_order 
	FROM workflow_items wi
	INNER JOIN workflows w ON wi.parent_id = w.id
	INNER JOIN documents d ON wi.document_id = d.id
	WHERE 
		w.module_id = ? AND w.store_id = ? AND 
		wi.document_id = ? AND wi.account_id = ? AND
		d.status = ?
	ORDER BY approval_order
	LIMIT 1
	`

	var appOrder int

	err := conn.QueryRow(query, s.GetModuleID(), storeID, docID, accountID, tr.Status()).Scan(&appOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("get approval order: [%w], query: %s", err, query)
	}

	logger.Debug("approval order", appOrder)

	if appOrder == 1 {
		return true, nil
	}

	return false, nil
}

func (s *workflowRepository) checkActionUpdated(conn *sql.DB, tr *treegrid.MainRow, accountID int) (bool, error) {
	query := `
	SELECT document_id, store_id, status 
	FROM transfers
	WHERE id = ?
	`

	var docID, storeID, status int
	if err := conn.QueryRow(query, tr.Fields.GetID()).Scan(&docID, &storeID, &status); err != nil {
		return false, fmt.Errorf("query row: [%w]", err)
	}

	if newDocID, ok := tr.Fields.GetValInt("document_id"); ok {
		docID = newDocID
	}

	if newStoreID, ok := tr.Fields.GetValInt("store_id"); ok {
		storeID = newStoreID
	}

	newStatus, ok := tr.Fields.GetValInt("status")
	if !ok {
		newStatus = status
	}
	log.Println("Prev status", status, "New status", newStatus)

	// status stays the same so need to check if user has permission to this status without approval order
	if newStatus == status {
		log.Println("Status stays the same")

		// when status 0 all operations are allowed
		if status == 0 {
			return true, nil
		}

		// only status equel to 1 to operation is allowed
		if status == 1 {
			return false, nil
		}

		if status > 1 {
			return false, nil
		}

		return true, nil
	}

	// no updates or adding new items are allowed
	if newStatus == 1 {
		if len(tr.Items) > 0 {
			logger.Debug("can't update or add items with status 1")

			return false, nil
		}

		updatedFields := tr.Fields.UpdatedFields()
		if len(updatedFields) != 1 {
			logger.Debug("can't do updates with status 1", updatedFields)

			return false, nil
		}

		if updatedFields[0] != "status" {
			logger.Debug("can update only status")

			return false, nil
		}

		return true, nil
	}

	logger.Debug("without statuses > 1")

	return false, nil
}
