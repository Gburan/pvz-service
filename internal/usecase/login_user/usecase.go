package login_user

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
	record, err := u.repoUser.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository2.ErrUserNotFound) {
			return nil, fmt.Errorf("%w: %s", usecase2.ErrNotFoundUser, req.Email)
		}
		return nil, fmt.Errorf("%w: %v", usecase2.ErrGetUser, err)
	}

	if err = compareHashPass(record.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("%w: %v", usecase2.ErrIncorrectPass, err)
	}

	return &Out{User: *record}, nil
}

func compareHashPass(dbhashedPass, reqPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbhashedPass), []byte(reqPass))
}
