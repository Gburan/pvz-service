package pvz

import (
	"context"

	"pvz-service/internal/model/entity"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=pvz RepositoryPVZ
type RepositoryPVZ interface {
	SavePVZ(ctx context.Context, city string) (*entity.PVZ, error)
	GetPVZByID(ctx context.Context, pvzId string) (*entity.PVZ, error)
	GetPVZsByIDs(ctx context.Context, pvzIds []string) (*[]entity.PVZ, error)
}
