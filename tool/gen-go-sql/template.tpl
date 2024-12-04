// Code generated by gen-go-datastore-sql. DO NOT EDIT.

package {{ .Package }}

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

type client interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type {{ .CamelTableName }} struct {
	dbClient client
	logger   *log.Logger
}

func New{{ .PascalTableName }}(dbClient client, logger *log.Logger) repository.{{ .PascalTableName }} {
	return &{{ .CamelTableName }}{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (m *{{ .CamelTableName }}) Get(ctx context.Context, id string, opt ...qm.QueryMod) (*model.{{ .PascalTableName }}, error) {
	query := make([]qm.QueryMod, 0, len(opt)+2)
	query = append(query, qm.Where("id = ?", id), qm.And("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will exec {{ .PascalTableName }}.Get, package: {{ .Package }}")

	{{ .CamelTableName }}, err := model.{{ .PluralPascalTableName }}(
		query...,
	).One(ctx, m.dbClient)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found {{ .CamelTableName }}: %w", err)
		}

		return nil, fmt.Errorf("error executing {{ .CamelTableName }}.Get: %w", err)
	}

	return {{ .CamelTableName }}, nil
}

func (u *{{ .CamelTableName }}) BatchGet(ctx context.Context, ids []any, opt ...qm.QueryMod) (model.{{ .PascalTableName }}Slice, error) {
	query := make([]qm.QueryMod, 0, len(opt)+2)
	query = append(query, qm.WhereIn("id IN ?", ids...), qm.And("deleted_at IS NULL"))
	query = append(query, opt...)

	u.logger.Printf("Will exec {{ .PascalTableName }}.BatchGet, package: {{ .Package }}")

	{{ .PluralCamelTableName }}, err := model.{{ .PluralPascalTableName }}(
		query...,
	).All(ctx, u.dbClient)
	if err != nil {
		return nil, fmt.Errorf("error executing {{ .CamelTableName }}.BatchGet: %w", err)
	}

	return {{ .PluralCamelTableName }}, nil
}

func (m *{{ .CamelTableName }}) Count(ctx context.Context, opt ...qm.QueryMod) (int64, error) {
	query := make([]qm.QueryMod, 0, len(opt)+1)
	query = append(query, qm.Where("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will execute {{ .PascalTableName }}.Count, package: {{ .Package }}")

	count, err := model.{{ .PluralPascalTableName }}(
		query...,
	).Count(ctx, m.dbClient)
	if err != nil {
		return 0, fmt.Errorf("error executing {{ .CamelTableName }}.Count: %w", err)
	}

	return count, nil
}

func (m *{{ .CamelTableName }}) List(ctx context.Context, opt ...qm.QueryMod) (model.{{ .PascalTableName }}Slice, error) {
	query := make([]qm.QueryMod, 0, len(opt)+1)
	query = append(query, qm.Where("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will execute {{ .PascalTableName }}.List, package: {{ .Package }}")

	{{ .PluralCamelTableName }}, err := model.{{ .PluralPascalTableName }}(
		query...,
	).All(ctx, m.dbClient)
	if err != nil {
		return nil, fmt.Errorf("error executing {{ .CamelTableName }}.List: %w", err)
	}

	return {{ .PluralCamelTableName }}, nil
}

func (m *{{ .CamelTableName }}) Insert(ctx context.Context, {{ .CamelTableName }} *model.{{ .PascalTableName }}) error {
	m.logger.Printf("Will execute {{ .PascalTableName }}.Insert, package: {{ .Package }}")

	err := {{ .CamelTableName }}.Insert(ctx, m.dbClient, boil.Infer())
	if err != nil {
		return fmt.Errorf("error executing {{ .CamelTableName }}.Insert: %w", err)
	}

	return nil
}

func (u *{{ .CamelTableName }}) Delete(ctx context.Context, id string) error {
	u.logger.Printf("Will execute {{ .PascalTableName }}.Delete, package: {{ .Package }}")

	const query = `UPDATE {{ .TableName }} SET deleted_at = $1 WHERE id = $2`

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
		return fmt.Errorf("error executing {{ .CamelTableName }}.Delete: %w", err)
	}

	return nil
}

func (u *{{ .CamelTableName }}) BulkDelete(ctx context.Context, ids []any) error {
	if len(ids) == 0 {
		return nil
	}

	u.logger.Printf("Will execute {{ .PascalTableName }}.BulkDelete, package: {{ .Package }}")

	const query = `UPDATE {{ .TableName }} SET deleted_at = $1 WHERE id IN (?)`

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
		return fmt.Errorf("error executing {{ .CamelTableName }}.BulkDelete: %w", err)
	}

	return nil
}

func (m *{{ .CamelTableName }}) Exist(ctx context.Context, opt ...qm.QueryMod) (bool, error) {
	query := make([]qm.QueryMod, 0, len(opt)+1)
	query = append(query, qm.Where("deleted_at IS NULL"))
	query = append(query, opt...)

	m.logger.Printf("Will execute {{ .PascalTableName }}.Exist, package: {{ .Package }}")

	exists, err := model.{{ .PluralPascalTableName }}(
		query...,
	).Exists(ctx, m.dbClient)
	if err != nil {
		return false, fmt.Errorf("error executing {{ .CamelTableName }}.Exist: %w", err)
	}

	return exists, nil
}

