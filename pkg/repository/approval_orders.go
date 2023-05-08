package repository

type IdentityStorage interface {
	GetStatus(id interface{}) (int, error)
	GetDocID(id interface{}) (int, error)
}
type ApprovalOrders struct {
	WorkflowRepository
	IdentityStorage
}

func NewApprovalOrder(wrk WorkflowRepository, i IdentityStorage) *ApprovalOrders {
	return &ApprovalOrders{wrk, i}
}
