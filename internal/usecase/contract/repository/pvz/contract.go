package pvz

import (
	"context"

	"pvz-service/internal/model/entity"

	"github.com/google/uuid"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=pvz RepositoryPVZ
type RepositoryPVZ interface {
	SavePVZ(ctx context.Context, pvz entity.PVZ) (*entity.PVZ, error)
	GetPVZByID(ctx context.Context, pvz entity.PVZ) (*entity.PVZ, error)
	GetPVZsByIDs(ctx context.Context, pvzIds []uuid.UUID) (*[]entity.PVZ, error)
	GetPVZList(ctx context.Context) ([]*entity.PVZ, error)
}
