package reception

import "time"

type receptionDB struct {
	Uuid     string    `db:"id"`
	DateTime time.Time `db:"date_time"`
	PVZID    string    `db:"pvz_id"`
	Status   string    `db:"status"`
}
