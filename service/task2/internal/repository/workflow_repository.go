package repository

type workflowRepository struct {
}

func NewWorkflowRepository() WorkflowRepository {
	return &workflowRepository{}
}

// CheckApprovalOrder implements WorkflowRepository
func (wr *workflowRepository) CheckApprovalOrder(accountID int, status int) (bool, error) {
	return true, nil
}

// GetModuleID implements WorkflowRepository
func (wr *workflowRepository) GetModuleID() (int, error) {
	return 6, nil
}
