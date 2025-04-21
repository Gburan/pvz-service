package pvz_info

import (
	"time"

	"pvz-service/internal/model/entity"
)

type In struct {
	StartData time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}

type Out struct {
	PVZ        entity.PVZ
	Receptions []entity.ReceptionWithProducts
}
