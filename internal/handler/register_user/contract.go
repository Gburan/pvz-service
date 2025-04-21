package register_user

import (
	"context"

	"pvz-service/internal/usecase/register_user"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=register_user usecase
type usecase interface {
	Run(ctx context.Context, req register_user.In) (*register_user.Out, error)
}

type registerUserIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof_user"`
}

type registerUserOut struct {
	Uuid  string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
