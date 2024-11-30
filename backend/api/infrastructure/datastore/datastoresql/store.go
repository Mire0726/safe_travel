package datastoresql

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Store) ReadWriteStore() datastore.ReadWriteStore {
	return &nonTransactionalReadWriteStore{
		c:      s.dbClient,
		logger: s.logger,
	}
}

func (s *Store) ReadWriteTransaction(ctx context.Context, f func(context.Context, datastore.ReadWriteStore) error) error {
	tx, err := s.dbClient.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}

	defer func() {
		if e := tx.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			s.logger.Printf("failed to rollback transaction: %v", e)
		}
	}()

	rw := &transactionalReadWriteStore{
		tx:     tx,
		logger: s.logger,
	}

	if err := f(ctx, rw); err != nil {
		s.logger.Printf("failed to execute transaction: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Printf("failed to commit transaction: %v", err)

		return err
	}

	return tx.Commit()
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

type transactionalReadWriteStore struct {
	tx     *sql.Tx
	logger *log.Logger
}

func (s *transactionalReadWriteStore) User() repository.User {
	return user.NewUser(s.tx, s.logger)
}

func (s *transactionalReadWriteStore) Event() repository.Event {
	return event.NewEvent(s.tx, s.logger)
}

func (s *transactionalReadWriteStore) Transport() repository.Transport {
	return transport.NewTransport(s.tx, s.logger)
}
