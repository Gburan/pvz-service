package login_user

import (
	"context"

	"pvz-service/internal/usecase/login_user"
)

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=login_user usecase
type usecase interface {
	Run(ctx context.Context, req login_user.In) (*login_user.Out, error)
}
