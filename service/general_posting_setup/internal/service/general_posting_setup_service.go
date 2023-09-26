package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type generalPostingSetupService struct {
	generalPostingSetupSimpleRepository treegrid.SimpleGridRowRepository
}

func NewGeneralPostingSetupService(generalPostingSetupSimpleRepository treegrid.SimpleGridRowRepository) GeneralPostingSetupService {
	return &generalPostingSetupService{
		generalPostingSetupSimpleRepository: generalPostingSetupSimpleRepository,
	}
}

// GetPageCount implements GeneralPostingSetupService
func (g *generalPostingSetupService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return g.generalPostingSetupSimpleRepository.GetPageCount(tr)
}

// GetPageData implements GeneralPostingSetupService
func (g *generalPostingSetupService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return g.generalPostingSetupSimpleRepository.GetPageData(tr)
}
