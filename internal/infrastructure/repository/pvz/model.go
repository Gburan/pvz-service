package pvz

import (
	"time"

	"github.com/google/uuid"
)

type pvzDB struct {
	Uuid             uuid.UUID `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
