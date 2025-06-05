package list_pvzs

import "pvz-service/internal/model/entity"

type Out struct {
	PVZs []*entity.PVZ
}
