package dummy_login

type dummyLoginIn struct {
	Role string `json:"role" validate:"required,oneof_user"`
}

type dummyLoginOut struct {
	Token string `json:"token"`
}
