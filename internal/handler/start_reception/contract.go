package start_reception

import (
	"context"

	"pvz-service/internal/usecase/start_reception"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=start_reception usecase
type usecase interface {
	Run(ctx context.Context, req start_reception.In) (*start_reception.Out, error)
}

type startReceptionIn struct {
	PVZID string `json:"pvzId" validate:"required,uuid4"`
}

type startReceptionOut struct {
	Id       string `json:"id"`
	DateTime string `json:"dateTime"`
	PvzId    string `json:"pvzId"`
	Status   string `json:"status"`
}
