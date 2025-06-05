package login_user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/logging"
	"pvz-service/internal/model/entity"
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
	slog.DebugContext(ctx, "Call GetUserByEmail")
	record, err := u.repoUser.GetUserByEmail(ctx, entity.User{
		Email: req.Email,
	})
	if err != nil {
		if errors.Is(err, repository2.ErrUserNotFound) {
			return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrNotFoundUser, req.Email))
		}
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrGetUser, err))
	}

	slog.DebugContext(ctx, "compareHashPass")
	if err = compareHashPass(record.PasswordHash, req.Password); err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrIncorrectPass, err))
	}

	slog.DebugContext(ctx, "Usecase Login user success")
	return &Out{User: *record}, nil
}

func compareHashPass(dbhashedPass, reqPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbhashedPass), []byte(reqPass))
}
