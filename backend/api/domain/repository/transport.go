package repository

import (
	"context"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Transport interface {
	Insert(ctx context.Context, m *model.Transport) error
	List(ctx context.Context, opt ...qm.QueryMod) (model.TransportSlice, error)
	Delete(ctx context.Context, id string) error
	Exist(ctx context.Context, opt ...qm.QueryMod) (bool, error)
}
