package delete_product

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

func (u *usecase) Run(ctx context.Context, req In) error {
	_, err := u.repoPVZ.GetPVZByID(ctx, req.PVZID)
	if err != nil {
		if errors.Is(err, repository2.ErrPVZNotFound) {
			return fmt.Errorf("%w: %s", usecase2.ErrNotFoundPVZ, req.PVZID)
		}
		return fmt.Errorf("%w: %s", usecase2.ErrGetPVZByID, req.PVZID)
	}

	lastReception, err := u.repoReception.GetLastReceptionPVZ(ctx, req.PVZID)
	if err != nil {
		if errors.Is(err, repository2.ErrReceptionNotFound) {
			return fmt.Errorf("%w: %s", usecase2.ErrNotFoundReception, req.PVZID)
		}
		return fmt.Errorf("%w: %s", usecase2.ErrGetReception, req.PVZID)
	}
	if lastReception.Status == statusReceptionDone {
		return fmt.Errorf("%w: %s", usecase2.ErrNotFoundOpenedReception, req.PVZID)
	}

	lastProduct, err := u.repoProduct.GetLastProductByReceptionPVZ(ctx, lastReception.Uuid)
	if err != nil {
		if errors.Is(err, repository2.ErrProductNotFound) {
			return fmt.Errorf("%w: %s", usecase2.ErrNotFoundProduct, req.PVZID)
		}
		return fmt.Errorf("%w: %s", usecase2.ErrGetProduct, req.PVZID)
	}

	err = u.repoProduct.DeleteProduct(ctx, lastProduct.Uuid)
	if err != nil {
		return fmt.Errorf("%w: %s", usecase2.ErrDeleteProduct, req.PVZID)
	}
	return nil
}
