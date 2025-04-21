package add_product

import (
	"context"

	"pvz-service/internal/usecase/add_product"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=add_product usecase
type usecase interface {
	Run(ctx context.Context, req add_product.In) (*add_product.Out, error)
}

type addProductIn struct {
	Type  string `json:"type" validate:"required,oneof_category"`
	PVZID string `json:"pvzId" validate:"required,uuid4"`
}

type addProductOut struct {
	Uuid        string `json:"id"`
	DateTime    string `json:"dateTime"`
	Type        string `json:"type"`
	ReceptionID string `json:"receptionId"`
}
