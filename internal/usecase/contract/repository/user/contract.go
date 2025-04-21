package user

import (
	"context"

	"pvz-service/internal/model/entity"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=user RepositoryReception
type RepositoryUser interface {
	AddUser(ctx context.Context, email string, passwordHash string, role string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
