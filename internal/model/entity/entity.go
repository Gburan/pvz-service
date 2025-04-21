package entity

import "time"

type User struct {
	ID           string
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
	Uuid             string
	RegistrationDate time.Time
	City             string
}

type Reception struct {
	Uuid     string
	DateTime time.Time
	PVZID    string
	Status   string
}

type Product struct {
	Uuid        string
	DateTime    time.Time
	Type        string
	ReceptionID string
}
