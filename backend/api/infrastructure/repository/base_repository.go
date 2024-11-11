// api/infrastructure/repository/base_repository.go
package repository

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// BaseModel はすべてのモデルの基本インターフェース
type BaseModel interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
	Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error)
	Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error)
}

// BaseRepository は基本的なCRUD操作を提供する
type BaseRepository struct {
	executor boil.ContextExecutor
}

func NewBaseRepository(executor boil.ContextExecutor) *BaseRepository {
	return &BaseRepository{
		executor: executor,
	}
}

func (r *BaseRepository) Create(ctx context.Context, model BaseModel) error {
	return model.Insert(ctx, r.executor, boil.Infer())
}

func (r *BaseRepository) Update(ctx context.Context, model BaseModel) error {
	_, err := model.Update(ctx, r.executor, boil.Infer())
	return err
}

func (r *BaseRepository) Delete(ctx context.Context, model BaseModel) error {
	_, err := model.Delete(ctx, r.executor)
	return err
}

// Transaction はトランザクションを扱うためのヘルパー関数
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	if db, ok := r.executor.(*sql.DB); ok {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if err := fn(tx); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return rbErr
			}
			return err
		}

		return tx.Commit()
	}
	return fn(r.executor.(*sql.Tx))
}
