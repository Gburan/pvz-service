package register_user

import (
	"context"
	"errors"
	"fmt"

	repository2 "pvz-service/internal/infrastructure/repository"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/contract/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type usecase struct {
	repoUser user.RepositoryUser
}

func NewUsecase(repoUser user.RepositoryUser) *usecase {
	return &usecase{
		repoUser: repoUser,
	}
}

func (u *usecase) Run(ctx context.Context, req In) (*Out, error) {
	exist, err := u.repoUser.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repository2.ErrUserNotFound) {
		return nil, fmt.Errorf("%w: %v", usecase2.ErrGetUser, err)
	}
	if exist != nil {
		return nil, fmt.Errorf("%w: %s", usecase2.ErrUserAlreadyExist, req.Email)
	}

	hashedPass, err := generateHashPass(req.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", usecase2.ErrGenHashedPass, err)
	}

	record, err := u.repoUser.AddUser(ctx, req.Email, hashedPass, req.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", usecase2.ErrAddUser, err)
	}

	return &Out{User: *record}, nil
}

func generateHashPass(reqPass string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(reqPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}
