package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	productsTableName     = "product"
	idColumnName          = "id"
	dateTimeColumnName    = "date_time"
	typeColumnName        = "type"
	receptionIdColumnName = "reception_id"
)

type repository struct {
	db repository2.DBContract
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}

func (r *repository) DeleteProduct(ctx context.Context, prodId string) error {
	queryBuilder := squirrel.Delete(productsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{idColumnName: prodId})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}

	return nil
}

func (r *repository) GetLastProductByReceptionPVZ(ctx context.Context, recId string) (*entity.Product, error) {
	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, typeColumnName, receptionIdColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(productsTableName).
		Where(squirrel.Eq{receptionIdColumnName: recId}).
		OrderBy(fmt.Sprintf("%s %s", dateTimeColumnName, "DESC")).
		Limit(1)

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[productDB])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrProductNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.Product{
		Uuid:        result.Uuid,
		DateTime:    result.DateTime,
		Type:        result.Type,
		ReceptionID: result.ReceptionID,
	}, nil
}

func (r *repository) AddProduct(ctx context.Context, reception, tpe string) (*entity.Product, error) {
	queryBuilder := squirrel.Insert(productsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(typeColumnName, receptionIdColumnName).
		Values(tpe, reception).
		Suffix(
			fmt.Sprintf("RETURNING %s, %s, %s, %s",
				idColumnName,
				dateTimeColumnName,
				typeColumnName,
				receptionIdColumnName,
			),
		)
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[productDB])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.Product{
		Uuid:        result.Uuid,
		DateTime:    result.DateTime,
		Type:        result.Type,
		ReceptionID: result.ReceptionID,
	}, nil
}

func (r *repository) GetProductsByTimeRange(ctx context.Context, startDate, endData time.Time) (*[]entity.Product, error) {
	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, typeColumnName, receptionIdColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(productsTableName).
		Where(squirrel.And{
			squirrel.GtOrEq{dateTimeColumnName: startDate},
			squirrel.LtOrEq{dateTimeColumnName: endData},
		})

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[productDB])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrProductsNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	returning := make([]entity.Product, 0, len(result))
	for _, item := range result {
		retItem := entity.Product{
			Uuid:        item.Uuid,
			DateTime:    item.DateTime,
			Type:        item.Type,
			ReceptionID: item.ReceptionID,
		}
		returning = append(returning, retItem)
	}
	return &returning, nil
}
