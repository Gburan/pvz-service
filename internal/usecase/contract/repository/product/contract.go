package product

import (
	"context"
	"time"

	"pvz-service/internal/model/entity"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=product RepositoryProduct
type RepositoryProduct interface {
	GetLastProductByReceptionPVZ(ctx context.Context, product entity.Product) (*entity.Product, error)
	AddProduct(ctx context.Context, product entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, product entity.Product) error
	GetProductsByTimeRange(ctx context.Context, startDate, endDate time.Time) (*[]entity.Product, error)
}
