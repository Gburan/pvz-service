package register_user

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
	exist, err := u.repoUser.GetUserByEmail(ctx, entity.User{
		Email: req.Email,
	})
	if err != nil && !errors.Is(err, repository2.ErrUserNotFound) {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrGetUser, err))
	}
	if exist != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrUserAlreadyExist, req.Email))
	}

	slog.DebugContext(ctx, "generateHashPass")
	hashedPass, err := generateHashPass(req.Password)
	if err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrGenHashedPass, err))
	}

	slog.DebugContext(ctx, "Call AddUser")
	record, err := u.repoUser.AddUser(ctx, entity.User{
		Email:        req.Email,
		PasswordHash: hashedPass,
		Role:         req.Role,
	})
	if err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrAddUser, err))
	}

	slog.DebugContext(ctx, "Usecase Register user success")
	return &Out{User: *record}, nil
}

func generateHashPass(reqPass string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(reqPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}
