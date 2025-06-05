package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTableName     = "pvzuser"
	idColumnName       = "id"
	emailColumnName    = "email"
	passHashColumnName = "pass_hash"
	roleColumnName     = "role"

	returnAll = "RETURNING *"
)

type repository struct {
	db repository2.DBContract
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}

func (r *repository) AddUser(ctx context.Context, user entity.User) (*entity.User, error) {
	if user.Uuid == uuid.Nil {
		user.Uuid = uuid.New()
	}

	queryBuilder := squirrel.Insert(usersTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(idColumnName, emailColumnName, passHashColumnName, roleColumnName).
		Values(user.Uuid, user.Email, user.PasswordHash, user.Role).
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

	_, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[userDB])
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository AddUser success")
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, user entity.User) (*entity.User, error) {
	selectBuilder := squirrel.
		Select(idColumnName, emailColumnName, passHashColumnName, roleColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(usersTableName).
		Where(squirrel.Eq{emailColumnName: user.Email})

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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userDB])
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	slog.DebugContext(ctx, "Repository GetUserByEmail success")
	return &entity.User{
		Uuid:         result.Uuid,
		Email:        result.Email,
		PasswordHash: result.PassHash,
		Role:         result.Role,
	}, nil
}
