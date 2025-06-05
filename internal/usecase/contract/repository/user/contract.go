package user

import (
	"context"

	"pvz-service/internal/model/entity"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=user RepositoryReception
type RepositoryUser interface {
	AddUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, user entity.User) (*entity.User, error)
}
