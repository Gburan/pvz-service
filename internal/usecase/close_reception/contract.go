package close_reception

import "pvz-service/internal/model/entity"

type In struct {
	PVZID string
}

type Out struct {
	Reception entity.Reception
}
