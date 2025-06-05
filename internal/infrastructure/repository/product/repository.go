package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	"pvz-service/internal/usecase/contract/nower"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	productsTableName     = "product"
	idColumnName          = "id"
	dateTimeColumnName    = "date_time"
	typeColumnName        = "type"
	receptionIdColumnName = "reception_id"

	returnAll = "RETURNING *"
)

type repository struct {
	db    repository2.DBContract
	nower nower.Nower
}

func NewRepository(pool *pgxpool.Pool, nower nower.Nower) *repository {
	return &repository{
		db:    pool,
		nower: nower,
	}
}

func (r *repository) DeleteProduct(ctx context.Context, product entity.Product) error {
	queryBuilder := squirrel.Delete(productsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{idColumnName: product.Uuid})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		return fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		return fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}

	slog.DebugContext(ctx, "Repository DeleteProduct success")
	return nil
}

func (r *repository) GetLastProductByReceptionPVZ(ctx context.Context, product entity.Product) (*entity.Product, error) {
	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, typeColumnName, receptionIdColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(productsTableName).
		Where(squirrel.Eq{receptionIdColumnName: product.ReceptionID}).
		OrderBy(fmt.Sprintf("%s %s", dateTimeColumnName, "DESC")).
		Limit(1)

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		slog.DebugContext(ctx, err.Error())

		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[productDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrProductNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository GetLastProductByReceptionPVZ success")
	return &entity.Product{
		Uuid:        result.Uuid,
		DateTime:    result.DateTime,
		Type:        result.Type,
		ReceptionID: result.ReceptionID,
	}, nil
}

func (r *repository) AddProduct(ctx context.Context, product entity.Product) (*entity.Product, error) {
	if product.Uuid == uuid.Nil {
		product.Uuid = uuid.New()
	}
	if product.DateTime.IsZero() {
		product.DateTime = r.nower.Now()
	}

	queryBuilder := squirrel.Insert(productsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(idColumnName, dateTimeColumnName, typeColumnName, receptionIdColumnName).
		Values(product.Uuid.String(), product.DateTime, product.Type, product.ReceptionID.String()).
		Suffix(returnAll)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	_, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[productDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository Addproduct success")
	return &product, nil
}

func (r *repository) GetProductsByTimeRange(ctx context.Context, startDate, endDate time.Time) (*[]entity.Product, error) {
	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, typeColumnName, receptionIdColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(productsTableName).
		Where(squirrel.And{
			squirrel.GtOrEq{dateTimeColumnName: startDate},
			squirrel.LtOrEq{dateTimeColumnName: endDate},
		})

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[productDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
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

	slog.DebugContext(ctx, "Repository GetProductsByTimeRange success")
	return &returning, nil
}
