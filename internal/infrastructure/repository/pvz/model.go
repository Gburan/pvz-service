package pvz

import "time"

type pvzDB struct {
	Uuid             string    `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
