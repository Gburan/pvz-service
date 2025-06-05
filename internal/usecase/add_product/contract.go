package add_product

import (
	"pvz-service/internal/model/entity"

	"github.com/google/uuid"
)

type In struct {
	Type  string
	PVZID uuid.UUID
}

type Out struct {
	Product entity.Product
}
