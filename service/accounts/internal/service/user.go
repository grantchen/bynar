package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repos sync.Map
}

func NewUserService() *UserService {
	return &UserService{repos: sync.Map{}}
}

func (s *UserService) repo(ctx context.Context) (treegrid.SimpleGridRowRepository, error) {
	claims, ok := ctx.Value("id_token").(middleware.IdTokenClaims)
	if !ok {
		return nil, errors.New("no claims in context")
	}
	connStr := os.Getenv(claims.TenantUuid) + claims.OrganizationUuid
	conn, ok := s.repos.Load(connStr)
	var repo treegrid.SimpleGridRowRepository
	if !ok {
		db, err := sql_db.InitializeConnection(connStr)
		logrus.Info("init db ", connStr)
		if err != nil {
			return nil, err
		}
		repo := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "users", map[string][]string{},
			100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
		s.repos.Store(connStr, repo)
	} else {
		logrus.Info("repo exist ", connStr)
		repo = conn.(treegrid.SimpleGridRowRepository)
	}
	return repo, nil
}

func (s *UserService) Handle(ctx context.Context, req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	for _, gr := range grList {
		// TODO:
		// if err := u.handle(gr); err != nil {
		// 	log.Println("Err", err)

		// 	resp.IO.Result = -1
		// 	resp.IO.Message += err.Error() + "\n"
		// 	resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
		// 	break
		// }
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}

	return resp, nil
}

func (s *UserService) GetPageCount(ctx context.Context, tr *treegrid.Treegrid) (float64, error) {
	repo, err := s.repo(ctx)
	if err != nil {
		return 0, err
	}
	count, err := repo.GetPageCount(tr)
	return float64(count), err
}

func (s *UserService) GetPageData(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error) {
	repo, err := s.repo(ctx)
	if err != nil {
		return nil, err
	}
	return repo.GetPageData(tr)
}
