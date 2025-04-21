package create_pvz

import (
	"context"

	"pvz-service/internal/usecase/create_pvz"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=create_pvz usecase
type usecase interface {
	Run(ctx context.Context, req create_pvz.In) (*create_pvz.Out, error)
}

type createPVZIn struct {
	City string `json:"city" validate:"required,oneof_city"`
}

type createPVZOut struct {
	Uuid             string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}
