package create_pvz

import (
	"context"
	"fmt"

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
	record, err := u.repoPVZ.SavePVZ(ctx, req.City)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", usecase2.ErrAddPVZ, err)
	}

	return &Out{PVZ: *record}, nil
}
