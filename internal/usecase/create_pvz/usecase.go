package create_pvz

import (
	"context"
	"fmt"
	"log/slog"

	"pvz-service/internal/logging"
	"pvz-service/internal/metrics"
	"pvz-service/internal/model/entity"
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

func (u *usecase) Run(ctx context.Context, req In) (*Out, error) {
	slog.DebugContext(ctx, "Call SavePVZ")
	record, err := u.repoPVZ.SavePVZ(ctx, entity.PVZ{
		City: req.City,
	})
	if err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrAddPVZ, err))
	}

	metrics.IncCreatedPVZ(record.City)
	slog.DebugContext(ctx, "Usecase CreatePVZ success")
	return &Out{PVZ: *record}, nil
}
