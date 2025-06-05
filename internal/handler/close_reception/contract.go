package close_reception

import (
	"context"

	"pvz-service/internal/usecase/close_reception"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=close_reception usecase
type usecase interface {
	Run(ctx context.Context, req close_reception.In) (*close_reception.Out, error)
}
