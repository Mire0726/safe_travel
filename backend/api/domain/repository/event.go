package repository

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
)

type Event interface {
	Insert(ctx context.Context, m *model.Event) error
}
