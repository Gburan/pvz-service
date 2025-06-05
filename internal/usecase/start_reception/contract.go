package start_reception

import (
	"pvz-service/internal/model/entity"

	"github.com/google/uuid"
)

type In struct {
	PVZID uuid.UUID
}

type Out struct {
	Reception entity.Reception
}
