package list_pvzs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/contract/repository/pvz"
)

type usecase struct {
	repoPVZ pvz.RepositoryPVZ
}

func NewUsecase(repoPVZ pvz.RepositoryPVZ) *usecase {
	return &usecase{
		repoPVZ: repoPVZ,
	}
}

func (u *usecase) Run(ctx context.Context) (*Out, error) {
	slog.DebugContext(ctx, "Call List PVZs")
	record, err := u.repoPVZ.GetPVZList(ctx)
	if err != nil && !errors.Is(err, repository.ErrPVZNotFound) {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrListPVZs, err))
	}

	slog.DebugContext(ctx, "Usecase List PVZs success")
	return &Out{PVZs: record}, nil
}
