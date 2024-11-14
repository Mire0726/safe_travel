package datastore

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/repository"
)

type Data interface {
	ReadWriteStore() ReadWriteStore

	ReadWriteTransaction(ctx context.Context, f func(context.Context, ReadWriteStore) error) error
}
type ReadWriteStore interface {
	User() repository.User
	Event() repository.Event
	Transport() repository.Transport
}
