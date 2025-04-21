package add_product

import "pvz-service/internal/model/entity"

type In struct {
	Type  string
	PVZID string
}

type Out struct {
	Product entity.Product
}
