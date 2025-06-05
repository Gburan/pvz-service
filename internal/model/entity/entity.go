package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid         uuid.UUID
	Email        string
	PasswordHash string
	Role         string
}

type PVZWithReceptions struct {
	PVZ        PVZ
	Receptions []ReceptionWithProducts
}

type ReceptionWithProducts struct {
	Reception Reception
	Products  []Product
}

type PVZ struct {
	Uuid             uuid.UUID
	RegistrationDate time.Time
	City             string
}

type Reception struct {
	Uuid     uuid.UUID
	DateTime time.Time
	PVZID    uuid.UUID
	Status   string
}

type Product struct {
	Uuid        uuid.UUID
	DateTime    time.Time
	Type        string
	ReceptionID uuid.UUID
}
