package register_user

import (
	"context"

	"pvz-service/internal/usecase/register_user"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=register_user usecase
type usecase interface {
	Run(ctx context.Context, req register_user.In) (*register_user.Out, error)
}
