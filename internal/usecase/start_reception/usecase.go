package start_reception

import (
	"context"
	"errors"
	"fmt"

	repository2 "pvz-service/internal/infrastructure/repository"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/contract/repository/product"
	"pvz-service/internal/usecase/contract/repository/pvz"
	"pvz-service/internal/usecase/contract/repository/reception"
)

const (
	statusReceptionDone = "close"
)

type usecase struct {
	repoPVZ       pvz.RepositoryPVZ
	repoReception reception.RepositoryReception
	repoProduct   product.RepositoryProduct
}

func NewUsecase(
	repoPVZ pvz.RepositoryPVZ,
	repoReception reception.RepositoryReception,
	repoProduct product.RepositoryProduct,
) *usecase {
	return &usecase{
		repoPVZ:       repoPVZ,
		repoReception: repoReception,
		repoProduct:   repoProduct,
	}
}

func (u *usecase) Run(ctx context.Context, req In) (*Out, error) {
	_, err := u.repoPVZ.GetPVZByID(ctx, req.PVZID)
	if err != nil {
		if errors.Is(err, repository2.ErrPVZNotFound) {
			return nil, fmt.Errorf("%w: %s", usecase2.ErrNotFoundPVZ, req.PVZID)
		}
		return nil, fmt.Errorf("%w: %s", usecase2.ErrGetPVZByID, req.PVZID)
	}

	lastReception, err := u.repoReception.GetLastReceptionPVZ(ctx, req.PVZID)
	if err != nil && !errors.Is(err, repository2.ErrReceptionNotFound) {
		return nil, fmt.Errorf("%w: %s", usecase2.ErrGetReception, req.PVZID)
	}
	if lastReception != nil && lastReception.Status != statusReceptionDone {
		return nil, fmt.Errorf("%w: %s", usecase2.ErrFoundOpenedReception, req.PVZID)
	}

	record, err := u.repoReception.StartReception(ctx, req.PVZID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", usecase2.ErrStartReception, req.PVZID)
	}
	return &Out{Reception: *record}, nil
}
