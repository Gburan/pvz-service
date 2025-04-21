package user

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
	usersTableName     = "pvzuser"
	idColumnName       = "id"
	emailColumnName    = "email"
	passHashColumnName = "pass_hash"
	roleColumnName     = "role"
)

type repository struct {
	db repository2.DBContract
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}

func (r *repository) AddUser(ctx context.Context, email string, passwordHash string, role string) (*entity.User, error) {
	queryBuilder := squirrel.Insert(usersTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(emailColumnName, passHashColumnName, roleColumnName).
		Values(email, passwordHash, role).
		Suffix(
			fmt.Sprintf("RETURNING %s, %s, %s, %s",
				idColumnName,
				emailColumnName,
				passHashColumnName,
				roleColumnName,
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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userDB])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.User{
		ID:           result.Uuid,
		Email:        result.Email,
		PasswordHash: result.PassHash,
		Role:         result.Role,
	}, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	selectBuilder := squirrel.
		Select(idColumnName, emailColumnName, passHashColumnName, roleColumnName).
		PlaceholderFormat(squirrel.Dollar).
		From(usersTableName).
		Where(squirrel.Eq{emailColumnName: email})

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrBuildQuery, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository2.ErrExecuteQuery, err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userDB])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", repository2.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", repository2.ErrScanResult, err)
	}

	return &entity.User{
		ID:           result.Uuid,
		Email:        result.Email,
		PasswordHash: result.PassHash,
		Role:         result.Role,
	}, nil
}
