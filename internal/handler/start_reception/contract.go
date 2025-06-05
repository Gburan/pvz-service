package start_reception

import (
	"context"

	"pvz-service/internal/usecase/start_reception"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=start_reception usecase
type usecase interface {
	Run(ctx context.Context, req start_reception.In) (*start_reception.Out, error)
}
