package models

type WorkflowItem struct {
	Id            int
	ParentID      int
	AccountID     int
	DocumentID    int
	ApprovalOrder int
	Status        int
}
