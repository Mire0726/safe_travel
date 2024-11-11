package repository

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
)

type User interface {
	Insert(ctx context.Context, m *model.User) error
}
