// Code generated by gen-go-datastore-sql. DO NOT EDIT.

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Mire0726/safe_travel/backend/api/domain/model"
	"github.com/Mire0726/safe_travel/backend/api/domain/repository"
)

type event struct {
	dbClient client
	logger   *log.Logger
}

func NewEvent(dbClient client, logger *log.Logger) repository.Event {
	return &event{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (m *event) Get(ctx context.Context, id string, opt ...qm.QueryMod) (*model.Event, error) {
	query := make([]qm.QueryMod, 0, len(opt)+2)
	query = append(query, qm.Where("id = ?", id), qm.And("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will exec Event.Get, package: sql")

	event, err := model.Events(
		query...,
	).One(ctx, m.dbClient)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found event: %w", err)
		}

		return nil, fmt.Errorf("error executing event.Get: %w", err)
	}

	return event, nil
}

func (u *event) BatchGet(ctx context.Context, ids []any, opt ...qm.QueryMod) (model.EventSlice, error) {
	query := make([]qm.QueryMod, 0, len(opt)+2)
	query = append(query, qm.WhereIn("id IN ?", ids...), qm.And("deleted_at IS NULL"))
	query = append(query, opt...)

	u.logger.Printf("Will exec Event.BatchGet, package: sql")

	events, err := model.Events(
		query...,
	).All(ctx, u.dbClient)
	if err != nil {
		return nil, fmt.Errorf("error executing event.BatchGet: %w", err)
	}

	return events, nil
}

func (m *event) Count(ctx context.Context, opt ...qm.QueryMod) (int64, error) {
	query := make([]qm.QueryMod, 0, len(opt)+1)
	query = append(query, qm.Where("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will execute Event.Count, package: sql")

	count, err := model.Events(
		query...,
	).Count(ctx, m.dbClient)
	if err != nil {
		return 0, fmt.Errorf("error executing event.Count: %w", err)
	}

	return count, nil
}

func (m *event) List(ctx context.Context, opt ...qm.QueryMod) (model.EventSlice, error) {
	query := make([]qm.QueryMod, 0, len(opt)+1)
	query = append(query, qm.Where("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will execute Event.List, package: sql")

	events, err := model.Events(
		query...,
	).All(ctx, m.dbClient)
	if err != nil {
		return nil, fmt.Errorf("error executing event.List: %w", err)
	}

	return events, nil
}

func (m *event) Insert(ctx context.Context, event *model.Event) error {
	m.logger.Printf("Will execute Event.Insert, package: sql")

	err := event.Insert(ctx, m.dbClient, boil.Infer())
	if err != nil {
		return fmt.Errorf("error executing event.Insert: %w", err)
	}

	return nil
}

func (u *event) Delete(ctx context.Context, id string) error {
	u.logger.Printf("Will execute Event.Delete, package: sql")

	const query = `UPDATE events SET deleted_at = $1 WHERE id = $2`

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
	}

	_, err := u.dbClient.ExecContext(
		ctx,
		query,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("error executing event.Delete: %w", err)
	}

	return nil
}

func (u *event) BulkDelete(ctx context.Context, ids []any) error {
	if len(ids) == 0 {
		return nil
	}

	u.logger.Printf("Will execute Event.BulkDelete, package: sql")

	const query = `UPDATE events SET deleted_at = $1 WHERE id IN (?)`

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
	}

	_, err := u.dbClient.ExecContext(
		ctx,
		query,
		time.Now(),
		ids,
	)
	if err != nil {
		return fmt.Errorf("error executing event.BulkDelete: %w", err)
	}

	return nil
}