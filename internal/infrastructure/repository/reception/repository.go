package reception

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
	receptionsTableName = "reception"
	idColumnName        = "id"
	dateTimeColumnName  = "date_time"
	pvzIDColumnName     = "pvz_id"
	statusColumnName    = "status"

	statusReceptionDone     = "close"
	statusReceptionProgress = "in_progress"

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

func (r *repository) StartReception(ctx context.Context, reception entity.Reception) (*entity.Reception, error) {
	if reception.Uuid == uuid.Nil {
		reception.Uuid = uuid.New()
	}
	if reception.DateTime.IsZero() {
		reception.DateTime = r.nower.Now()
	}
	reception.Status = statusReceptionProgress

	queryBuilder := squirrel.Insert(receptionsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(idColumnName, dateTimeColumnName, pvzIDColumnName, statusColumnName).
		Values(reception.Uuid, reception.DateTime, reception.PVZID, reception.Status).
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

	_, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository StartReception success")
	return &reception, nil
}

func (r *repository) CloseReception(ctx context.Context, reception entity.Reception) (*entity.Reception, error) {
	queryBuilder := squirrel.Update(receptionsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Set(statusColumnName, statusReceptionDone).
		Where(squirrel.Eq{idColumnName: reception.Uuid}).
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
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	queryRes, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository CloseReception success")
	return &entity.Reception{
		Uuid:     queryRes.Uuid,
		DateTime: queryRes.DateTime,
		PVZID:    queryRes.PVZID,
		Status:   queryRes.Status,
	}, nil
}

func (r *repository) GetLastReceptionPVZ(ctx context.Context, reception entity.Reception) (*entity.Reception, error) {
	selectBuilder := squirrel.
		Select(idColumnName, dateTimeColumnName, pvzIDColumnName, statusColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(receptionsTableName).
		Where(squirrel.Eq{pvzIDColumnName: reception.PVZID}).
		OrderBy(fmt.Sprintf("%s %s", dateTimeColumnName, "DESC")).
		Limit(1)

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

	ret, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrReceptionNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository GetLastReceptionPVZ success")
	return &entity.Reception{
		Uuid:     ret.Uuid,
		DateTime: ret.DateTime,
		PVZID:    ret.PVZID,
		Status:   ret.Status,
	}, nil
}

func (r *repository) GetReceptionsByIDs(ctx context.Context, receptionIDs []uuid.UUID) (*[]entity.Reception, error) {
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
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	receptionsDB, err := pgx.CollectRows(rows, pgx.RowToStructByName[receptionDB])
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
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

	slog.DebugContext(ctx, "Repository GetReceptionsByIDs success")
	return &receptions, nil
}
