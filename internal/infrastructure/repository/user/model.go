package user

type userDB struct {
	Uuid     string `db:"id"`
	Email    string `db:"email"`
	PassHash string `db:"pass_hash"`
	Role     string `db:"role"`
}
