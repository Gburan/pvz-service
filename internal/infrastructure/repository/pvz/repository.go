package pvz

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	"pvz-service/internal/usecase/contract/nower"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pvzTableName               = "pvz"
	idColumnName               = "id"
	registrationDateColumnName = "registration_date"
	cityColumnName             = "city"

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

func (r *repository) SavePVZ(ctx context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
	if pvz.Uuid == uuid.Nil {
		pvz.Uuid = uuid.New()
	}
	if pvz.RegistrationDate.IsZero() {
		pvz.RegistrationDate = r.nower.Now()
	}

	queryBuilder := squirrel.Insert(pvzTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(idColumnName, registrationDateColumnName, cityColumnName).
		Values(pvz.Uuid, pvz.RegistrationDate, pvz.City).
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

	_, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository SavePVZ success")
	return &pvz, nil
}

func (r *repository) GetPVZByID(ctx context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
	selectBuilder := squirrel.
		Select(idColumnName, registrationDateColumnName, cityColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(pvzTableName).
		Where(squirrel.Eq{idColumnName: pvz.Uuid})

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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrPVZNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository GetPVZByID success")
	return &entity.PVZ{
		Uuid:             result.Uuid,
		RegistrationDate: result.RegistrationDate,
		City:             result.City,
	}, nil
}

func (r *repository) GetPVZsByIDs(ctx context.Context, pvzIds []uuid.UUID) (*[]entity.PVZ, error) {
	if len(pvzIds) == 0 {
		return &[]entity.PVZ{}, nil
	}

	selectBuilder := squirrel.
		Select(idColumnName, registrationDateColumnName, cityColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(pvzTableName).
		Where(squirrel.Eq{idColumnName: pvzIds})

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

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	pvzs := make([]entity.PVZ, 0, len(results))
	for _, result := range results {
		pvzs = append(pvzs, entity.PVZ{
			Uuid:             result.Uuid,
			RegistrationDate: result.RegistrationDate,
			City:             result.City,
		})
	}

	slog.DebugContext(ctx, "Repository GetPVZsByIDs success")
	return &pvzs, nil
}

func (r *repository) GetPVZList(ctx context.Context) ([]*entity.PVZ, error) {
	selectBuilder := squirrel.
		Select(idColumnName, registrationDateColumnName, cityColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(pvzTableName)

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

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrPVZNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	pvzs := make([]*entity.PVZ, 0, len(results))
	for _, result := range results {
		pvzs = append(pvzs, &entity.PVZ{
			Uuid:             result.Uuid,
			RegistrationDate: result.RegistrationDate,
			City:             result.City,
		})
	}

	slog.DebugContext(ctx, "Repository GetPVZList success")
	return pvzs, nil
}
