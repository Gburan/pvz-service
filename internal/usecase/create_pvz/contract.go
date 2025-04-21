package create_pvz

import (
	"pvz-service/internal/model/entity"
)

type In struct {
	City string
}

type Out struct {
	PVZ entity.PVZ
}
