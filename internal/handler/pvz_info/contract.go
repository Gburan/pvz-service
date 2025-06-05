package pvz_info

import (
	"context"

	"pvz-service/internal/usecase/pvz_info"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=pvz_info usecase
type usecase interface {
	Run(ctx context.Context, req pvz_info.In) ([]pvz_info.Out, error)
}
