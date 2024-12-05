package repository

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type User interface {
	Insert(ctx context.Context, m *model.User) error
	Get(ctx context.Context, id string, opt ...qm.QueryMod) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Exist(ctx context.Context, opt ...qm.QueryMod) (bool, error)
	Delete(ctx context.Context, id string) error
}
