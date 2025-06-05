package product

import (
	"time"

	"github.com/google/uuid"
)

type productDB struct {
	Uuid        uuid.UUID `db:"id"`
	DateTime    time.Time `db:"date_time"`
	Type        string    `db:"type"`
	ReceptionID uuid.UUID `db:"reception_id"`
}
