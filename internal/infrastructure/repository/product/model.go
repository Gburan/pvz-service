package product

import "time"

type productDB struct {
	Uuid        string    `db:"id"`
	DateTime    time.Time `db:"date_time"`
	Type        string    `db:"type"`
	ReceptionID string    `db:"reception_id"`
}
