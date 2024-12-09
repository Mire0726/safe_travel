package repository

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Event interface {
	Insert(ctx context.Context, m *model.Event) error
	List(ctx context.Context, opt ...qm.QueryMod) (model.EventSlice, error)
	Delete(ctx context.Context, id string) error
	Exist(ctx context.Context, opt ...qm.QueryMod) (bool, error)
}
