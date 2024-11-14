package datastoresql

import (
	"context"
	"database/sql"
	"log"

	"github.com/Mire0726/safe_travel/backend/api/domain/repository"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore/datastoresql/event"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore/datastoresql/transport"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore/datastoresql/user"
)

type client interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Store struct {
	dbClient client
	store    datastore.ReadWriteStore
	logger   *log.Logger
}

func NewStore(dbClient client, logger *log.Logger) *Store {
	return &Store{
		dbClient: dbClient,
		store: &nonTransactionalReadWriteStore{
			c:      dbClient,
			logger: logger,
		},
		logger: logger,
	}
}

func (s *Store) ReadWrite() datastore.ReadWriteStore {
	return s.store
}

type nonTransactionalReadWriteStore struct {
	c      client
	logger *log.Logger
}

func (s *nonTransactionalReadWriteStore) User() repository.User {
	return user.NewUser(s.c, s.logger)
}

func (s *nonTransactionalReadWriteStore) Event() repository.Event {
	return event.NewEvent(s.c, s.logger)
}

func (s *nonTransactionalReadWriteStore) Transport() repository.Transport {
	return transport.NewTransport(s.c, s.logger)
}
