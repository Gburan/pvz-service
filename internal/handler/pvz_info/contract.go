package pvz_info

import (
	"context"
	"time"

	"pvz-service/internal/usecase/pvz_info"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=pvz_info usecase
type usecase interface {
	Run(ctx context.Context, req pvz_info.In) ([]pvz_info.Out, error)
}

type pvzInfoIn struct {
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required"`
	Page      int       `json:"page" validate:"required"`
	Limit     int       `json:"limit" validate:"required"`
}

type pvzInfoOut struct {
	Pvz        pvzOut                     `json:"pvz"`
	Receptions []receptionWithProductsOut `json:"receptions"`
}

type receptionWithProductsOut struct {
	Reception receptionOut `json:"reception"`
	Products  []productOut `json:"products"`
}

type productOut struct {
	Uuid        string `json:"id"`
	DateTime    string `json:"dateTime"`
	Type        string `json:"type"`
	ReceptionID string `json:"receptionId"`
}

type receptionOut struct {
	Id       string `json:"id"`
	DateTime string `json:"dateTime"`
	PvzId    string `json:"pvzId"`
	Status   string `json:"status"`
}

type pvzOut struct {
	Uuid             string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}
