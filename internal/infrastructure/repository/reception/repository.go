package reception

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
	receptionsTableName = "reception"
	idColumnName        = "id"
	dateTimeColumnName  = "date_time"
	pvzIDColumnName     = "pvz_id"
	statusColumnName    = "status"

	statusReceptionDone = "close"
)

type repository struct {
	db repository2.DBContract
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}

func (r *repository) StartReception(ctx context.Context, pvzId string) (*entity.Reception, error) {
	queryBuilder := squirrel.Insert(receptionsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(pvzIDColumnName).
		Values(pvzId).
		Suffix(
			fmt.Sprintf("RETURNING %s, %s, %s, %s",
				idColumnName,
				dateTimeColumnName,
				pvzIDColumnName,
				statusColumnName,
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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.Reception{
		Uuid:     result.Uuid,
		DateTime: result.DateTime,
		PVZID:    result.PVZID,
		Status:   result.Status,
	}, nil
}

func (r *repository) CloseReception(ctx context.Context, recId string) (*entity.Reception, error) {
	queryBuilder := squirrel.Update(receptionsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Set(statusColumnName, statusReceptionDone).
		Where(squirrel.Eq{idColumnName: recId}).
		Suffix(
			fmt.Sprintf("RETURNING %s, %s, %s, %s",
				idColumnName,
				dateTimeColumnName,
				pvzIDColumnName,
				statusColumnName,
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

	queryRes, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.Reception{
		Uuid:     queryRes.Uuid,
		DateTime: queryRes.DateTime,
		PVZID:    queryRes.PVZID,
		Status:   queryRes.Status,
	}, nil
}

func (r *repository) GetLastReceptionPVZ(ctx context.Context, pvzId string) (*entity.Reception, error) {
	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, pvzIDColumnName, statusColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(receptionsTableName).
		Where(squirrel.Eq{pvzIDColumnName: pvzId}).
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

	reception, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrReceptionNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.Reception{
		Uuid:     reception.Uuid,
		DateTime: reception.DateTime,
		PVZID:    reception.PVZID,
		Status:   reception.Status,
	}, nil
}

func (r *repository) GetReceptionsByIDs(ctx context.Context, receptionIDs []string) (*[]entity.Reception, error) {
	if len(receptionIDs) == 0 {
		return &[]entity.Reception{}, nil
	}

	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, pvzIDColumnName, statusColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(receptionsTableName).
		Where(squirrel.Eq{idColumnName: receptionIDs})

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	receptionsDB, err := pgx.CollectRows(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	if len(receptionsDB) == 0 {
		return nil, fmt.Errorf("%w", repository2.ErrReceptionsNotFound)
	}

	receptions := make([]entity.Reception, 0, len(receptionsDB))
	for _, reception := range receptionsDB {
		receptions = append(receptions, entity.Reception{
			Uuid:     reception.Uuid,
			DateTime: reception.DateTime,
			PVZID:    reception.PVZID,
			Status:   reception.Status,
		})
	}
	return &receptions, nil
}
