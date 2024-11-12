package data

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/repository"
)

type Data interface {
	ReadWrite() ReadWrite

	ReadWriteTransaction(ctx context.Context, f func(context.Context, ReadWrite) error) error
}
type ReadWrite interface {
	User() repository.User
	Event() repository.Event
	Transport() repository.Transport
}
