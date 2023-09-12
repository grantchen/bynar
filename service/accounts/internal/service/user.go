package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/repository"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repos sync.Map
}

type Store struct {
	db   *sql.DB
	repo treegrid.SimpleGridRowRepository
}

func NewUserService() *UserService {
	return &UserService{repos: sync.Map{}}
}

func (s *UserService) store(ctx context.Context) (*Store, error) {
	claims, ok := ctx.Value("id_token").(middleware.IdTokenClaims)
	if !ok {
		return nil, errors.New("no claims in context")
	}
	connStr := os.Getenv(claims.TenantUuid) + claims.OrganizationUuid
	conn, ok := s.repos.Load(connStr)
	store := &Store{}
	var err error
	if !ok {
		store.db, err = sql_db.InitializeConnection(connStr)
		logrus.Info("init db ", connStr)
		if err != nil {
			return nil, err
		}
		store.repo = treegrid.NewSimpleGridRowRepositoryWithCfg(store.db, "users", map[string][]string{},
			100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
		s.repos.Store(connStr, store)
	} else {
		logrus.Info("repo exist ", connStr)
		store = conn.(*Store)
	}
	return store, nil
}

func (s *UserService) Handle(ctx context.Context, req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	grList, err := treegrid.ParseRequestUploadSingleRow(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	store, err := s.store(ctx)
	if err != nil {
		return nil, err
	}

	for _, gr := range grList {
		if err := s.handle(store, gr); err != nil {
			logrus.Error("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(gr))
			break
		}
		resp.Changes = append(resp.Changes, gr)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(gr))
	}

	return resp, nil
}

func (s *UserService) GetPageCount(ctx context.Context, tr *treegrid.Treegrid) (float64, error) {
	store, err := s.store(ctx)
	if err != nil {
		return 0, err
	}
	count, err := store.repo.GetPageCount(tr)
	return float64(count), err
}

func (s *UserService) GetPageData(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error) {
	store, err := s.store(ctx)
	if err != nil {
		return nil, err
	}
	return store.repo.GetPageData(tr)
}

func (s *UserService) handle(store *Store, gr treegrid.GridRow) error {
	tx, err := store.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	fieldsValidating := []string{"code"}

	// add addition here
	switch gr.GetActionType() {
	case treegrid.GridRowActionAdd:
		err1 := gr.ValidateOnRequiredAll(repository.UserFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := store.repo.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%v], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = store.repo.Add(tx, gr)
	case treegrid.GridRowActionChanged:
		err1 := gr.ValidateOnRequired(repository.UserFieldNames)
		if err1 != nil {
			return err1
		}
		ok, err1 := store.repo.ValidateOnIntegrity(gr, fieldsValidating)
		if !ok || err1 != nil {
			return fmt.Errorf("validate duplicate: [%w], field: %s", err1, strings.Join(fieldsValidating, ", "))
		}
		err = store.repo.Update(tx, gr)
	case treegrid.GridRowActionDeleted:
		err = store.repo.Delete(tx, gr)

	default:
		return fmt.Errorf("undefined row type: %s", gr.GetActionType())
	}

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: [%w]", err)
	}

	return err
}
