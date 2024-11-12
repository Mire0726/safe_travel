package repository

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
)

type Transport interface {
	Insert(ctx context.Context, m *model.Transport) error
}
