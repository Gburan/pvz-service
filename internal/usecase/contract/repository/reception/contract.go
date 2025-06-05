package reception

import (
	"context"

	"pvz-service/internal/model/entity"

	"github.com/google/uuid"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=reception RepositoryReception
type RepositoryReception interface {
	StartReception(ctx context.Context, reception entity.Reception) (*entity.Reception, error)
	CloseReception(ctx context.Context, reception entity.Reception) (*entity.Reception, error)
	GetLastReceptionPVZ(ctx context.Context, reception entity.Reception) (*entity.Reception, error)
	GetReceptionsByIDs(ctx context.Context, receptionIDs []uuid.UUID) (*[]entity.Reception, error)
}
