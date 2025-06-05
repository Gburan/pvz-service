package user

import "github.com/google/uuid"

type userDB struct {
	Uuid     uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	PassHash string    `db:"pass_hash"`
	Role     string    `db:"role"`
}
