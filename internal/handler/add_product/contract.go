package add_product

import (
	"context"

	"pvz-service/internal/usecase/add_product"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=add_product usecase
type usecase interface {
	Run(ctx context.Context, req add_product.In) (*add_product.Out, error)
}
