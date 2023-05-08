package sql_connection

import (
	"context"
	"database/sql"
	"sync"

	db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
)

type Pool struct {
	connections map[string]*sql.DB
	locker      *sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{locker: &sync.RWMutex{}, connections: map[string]*sql.DB{}}
}

type ErrCloseFailure struct {
	Errors []error
}

func (e ErrCloseFailure) Error() string {
	return "cloud not close all connections"
}

func (p *Pool) GetWithContext(_ context.Context, connectionString string) (*sql.DB, error) {
	p.locker.RLock()
	if v, exists := p.connections[connectionString]; exists {
		p.locker.RUnlock()

		return v, nil
	}

	p.locker.RUnlock()

	p.locker.Lock()
	defer p.locker.Unlock()

	connection, err := db.InitializeConnection(connectionString)
	if err != nil {
		return nil, err
	}

	p.connections[connectionString] = connection

	return connection, nil
}

func (p *Pool) Get(connectionString string) (*sql.DB, error) {
	p.locker.RLock()
	if v, exists := p.connections[connectionString]; exists {
		p.locker.RUnlock()

		return v, nil
	}

	p.locker.RUnlock()

	p.locker.Lock()
	defer p.locker.Unlock()

	connection, err := db.InitializeConnection(connectionString)
	if err != nil {
		return nil, err
	}

	p.connections[connectionString] = connection

	return connection, nil
}

func (p *Pool) Close() error {
	p.locker.Lock()
	defer p.locker.Unlock()

	var errs []error

	for _, connection := range p.connections {
		if closeErr := connection.Close(); closeErr != nil {
			errs = append(errs, closeErr)
		}
	}

	if len(errs) > 0 {
		return &ErrCloseFailure{Errors: errs}
	}

	p.connections = map[string]*sql.DB{}

	return nil
}
