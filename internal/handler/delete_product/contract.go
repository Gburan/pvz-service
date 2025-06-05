package delete_product

import (
	"context"

	"pvz-service/internal/usecase/delete_product"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=delete_product usecase
type usecase interface {
	Run(ctx context.Context, req delete_product.In) error
}
