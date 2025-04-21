package close_reception

import (
	"context"

	"pvz-service/internal/usecase/close_reception"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=close_reception usecase
type usecase interface {
	Run(ctx context.Context, req close_reception.In) (*close_reception.Out, error)
}

type closeReceptionIn struct {
	PVZID string `json:"pvzId" validate:"required,uuid4"`
}

type closeReceptionOut struct {
	Uuid     string `json:"id"`
	DateTime string `json:"dateTime"`
	PVZID    string `json:"pvzId"`
	Status   string `json:"status"`
}
