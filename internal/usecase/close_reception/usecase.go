package close_reception

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/logging"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/contract/repository/pvz"
	"pvz-service/internal/usecase/contract/repository/reception"
)

const (
	statusReceptionDone = "close"
)

type usecase struct {
	repoPVZ       pvz.RepositoryPVZ
	repoReception reception.RepositoryReception
}

func NewUsecase(
	repoPVZ pvz.RepositoryPVZ,
	repoReception reception.RepositoryReception,
) *usecase {
	return &usecase{
		repoPVZ:       repoPVZ,
		repoReception: repoReception,
	}
}

func (u *usecase) Run(ctx context.Context, req In) (*Out, error) {
	slog.DebugContext(ctx, "Call GetPVZByID")
	_, err := u.repoPVZ.GetPVZByID(ctx, entity.PVZ{
		Uuid: req.PVZID,
	})
	if err != nil {
		if errors.Is(err, repository2.ErrPVZNotFound) {
			return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrNotFoundPVZ, req.PVZID))
		}
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrGetPVZByID, req.PVZID))
	}

	slog.DebugContext(ctx, "Call GetLastReceptionPVZ")
	lastReception, err := u.repoReception.GetLastReceptionPVZ(ctx, entity.Reception{
		PVZID: req.PVZID,
	})
	if err != nil {
		if errors.Is(err, repository2.ErrReceptionNotFound) {
			return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrNotFoundReception, req.PVZID))
		}
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrGetReception, req.PVZID))
	}
	if lastReception.Status == statusReceptionDone {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrNotFoundOpenedReception, req.PVZID))
	}

	slog.DebugContext(ctx, "Call CloseReception")
	record, err := u.repoReception.CloseReception(ctx, entity.Reception{
		Uuid: lastReception.Uuid,
	})
	if err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %s", usecase2.ErrCloseReception, req.PVZID))
	}

	slog.DebugContext(ctx, "Usecase Close reception success")
	return &Out{Reception: *record}, nil
}
