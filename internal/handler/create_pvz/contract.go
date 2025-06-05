package create_pvz

import (
	"context"

	"pvz-service/internal/usecase/create_pvz"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=create_pvz usecase
type usecase interface {
	Run(ctx context.Context, req create_pvz.In) (*create_pvz.Out, error)
}
