package reception

import (
	"time"

	"github.com/google/uuid"
)

type receptionDB struct {
	Uuid     uuid.UUID `db:"id"`
	DateTime time.Time `db:"date_time"`
	PVZID    uuid.UUID `db:"pvz_id"`
	Status   string    `db:"status"`
}
