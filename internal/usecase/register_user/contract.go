package register_user

import "pvz-service/internal/model/entity"

type In struct {
	Email    string
	Password string
	Role     string
}

type Out struct {
	User entity.User
}
