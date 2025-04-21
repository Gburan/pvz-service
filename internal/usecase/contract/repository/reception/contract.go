package reception

import (
	"context"

	"pvz-service/internal/model/entity"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=reception RepositoryReception
type RepositoryReception interface {
	StartReception(ctx context.Context, pvzId string) (*entity.Reception, error)
	CloseReception(ctx context.Context, recId string) (*entity.Reception, error)
	GetLastReceptionPVZ(ctx context.Context, pvzId string) (*entity.Reception, error)
	GetReceptionsByIDs(ctx context.Context, receptionIDs []string) (*[]entity.Reception, error)
}
