package delete_product

import (
	"context"

	"pvz-service/internal/usecase/delete_product"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract_mock.go -package=delete_product usecase
type usecase interface {
	Run(ctx context.Context, req delete_product.In) error
}

type deleteProductIn struct {
	PVZID string `json:"pvzId" validate:"required,uuid4"`
}
