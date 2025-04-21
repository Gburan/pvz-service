package pvz

import (
	"context"
	"errors"
	"fmt"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pvzTableName               = "pvz"
	idColumnName               = "id"
	registrationDateColumnName = "registration_date"
	cityColumnName             = "city"
)

type repository struct {
	db repository2.DBContract
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}

func (r *repository) SavePVZ(ctx context.Context, city string) (*entity.PVZ, error) {
	queryBuilder := squirrel.Insert(pvzTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(cityColumnName).
		Values(city).
		Suffix(
			fmt.Sprintf("RETURNING %s, %s, %s",
				idColumnName,
				registrationDateColumnName,
				cityColumnName,
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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.PVZ{
		Uuid:             result.Uuid,
		RegistrationDate: result.RegistrationDate,
		City:             result.City,
	}, nil
}

func (r *repository) GetPVZByID(ctx context.Context, pvzId string) (*entity.PVZ, error) {
	selectBuilder := squirrel.
		Select(idColumnName, registrationDateColumnName, cityColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(pvzTableName).
		Where(squirrel.Eq{idColumnName: pvzId})

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrPVZNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.PVZ{
		Uuid:             result.Uuid,
		RegistrationDate: result.RegistrationDate,
		City:             result.City,
	}, nil
}

func (r *repository) GetPVZsByIDs(ctx context.Context, pvzIds []string) (*[]entity.PVZ, error) {
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
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[pvzDB])
	if err != nil {
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

	return &pvzs, nil
}
