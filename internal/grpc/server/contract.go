package server

import (
	"context"

	"pvz-service/internal/usecase/list_pvzs"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=server usecase
type usecase interface {
	Run(ctx context.Context) (*list_pvzs.Out, error)
}
