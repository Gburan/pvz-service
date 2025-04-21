package login_user

import (
	"context"

	"pvz-service/internal/usecase/login_user"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=login_user usecase
type usecase interface {
	Run(ctx context.Context, req login_user.In) (*login_user.Out, error)
}

type loginUserIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginUserOut struct {
	Token string `json:"token"`
}
