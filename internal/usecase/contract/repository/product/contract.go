package product

import (
	"context"
	"time"

	"pvz-service/internal/model/entity"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=product RepositoryProduct
type RepositoryProduct interface {
	GetLastProductByReceptionPVZ(ctx context.Context, recId string) (*entity.Product, error)
	AddProduct(ctx context.Context, reception, tpe string) (*entity.Product, error)
	DeleteProduct(ctx context.Context, prodId string) error
	GetProductsByTimeRange(ctx context.Context, startDate, endData time.Time) (*[]entity.Product, error)
}
