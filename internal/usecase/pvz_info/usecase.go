package pvz_info

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/logging"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/contract/repository/product"
	"pvz-service/internal/usecase/contract/repository/pvz"
	"pvz-service/internal/usecase/contract/repository/reception"

	"github.com/google/uuid"
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

func (u *usecase) Run(ctx context.Context, req In) ([]Out, error) {
	offset := calcOffset(req.Page, req.Limit)

	slog.DebugContext(ctx, "Call GetProductsByTimeRange")
	products, err := u.repoProduct.GetProductsByTimeRange(ctx, req.StartData, req.EndDate)
	if err != nil {
		if errors.Is(err, repository2.ErrProductsNotFound) {
			return nil, logging.WrapError(ctx, fmt.Errorf("%w between %v and %v", usecase2.ErrNotFoundProducts, req.StartData, req.EndDate))
		}
		return nil, logging.WrapError(ctx, fmt.Errorf("%ws between %v and %v", usecase2.ErrGetProducts, req.StartData, req.EndDate))
	}

	uniqueReceptionIds := make(map[uuid.UUID]struct{})
	for _, prod := range *products {
		uniqueReceptionIds[prod.ReceptionID] = struct{}{}
	}
	sequenceReceptionIds := make([]uuid.UUID, 0, len(uniqueReceptionIds))
	for id := range uniqueReceptionIds {
		sequenceReceptionIds = append(sequenceReceptionIds, id)
	}

	slog.DebugContext(ctx, "Call GetReceptionsByIDs")
	receptions, err := u.repoReception.GetReceptionsByIDs(ctx, sequenceReceptionIds)
	if err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrGetReceptions, err))
	}
	uniquePVZIds := make(map[uuid.UUID]struct{})
	receptionsMap := make(map[uuid.UUID]entity.Reception)
	pvzsIdsReceptions := make(map[uuid.UUID][]entity.Reception)
	for _, rec := range *receptions {
		receptionsMap[rec.Uuid] = rec
		pvzsIdsReceptions[rec.PVZID] = append(pvzsIdsReceptions[rec.PVZID], rec)
		uniquePVZIds[rec.PVZID] = struct{}{}
	}
	sequencePvzIds := make([]uuid.UUID, 0, len(uniquePVZIds))
	for id := range uniquePVZIds {
		sequencePvzIds = append(sequencePvzIds, id)
	}

	slog.DebugContext(ctx, "Call GetPVZsByIDs")
	pvzs, err := u.repoPVZ.GetPVZsByIDs(ctx, sequencePvzIds)
	if err != nil {
		return nil, logging.WrapError(ctx, fmt.Errorf("%w: %v", usecase2.ErrGetPVZs, err))
	}

	productsByReception := make(map[uuid.UUID][]entity.Product)
	for _, prod := range *products {
		productsByReception[prod.ReceptionID] = append(productsByReception[prod.ReceptionID], prod)
	}

	pvzSlice := *pvzs
	start := offset
	if start > len(*pvzs) {
		return nil, usecase2.ErrNotPageTooBig
	}
	end := start + req.Limit
	if end > len(*pvzs) {
		end = len(*pvzs)
	}

	var result []entity.PVZWithReceptions
	for _, pvz_ := range pvzSlice[start:end] {
		var receptionWithProducts []entity.ReceptionWithProducts

		for _, rec := range pvzsIdsReceptions[pvz_.Uuid] {
			if prods, exists := productsByReception[rec.Uuid]; exists {
				receptionWithProducts = append(receptionWithProducts, entity.ReceptionWithProducts{
					Reception: rec,
					Products:  prods,
				})
			}
		}

		result = append(result, entity.PVZWithReceptions{
			PVZ:        pvz_,
			Receptions: receptionWithProducts,
		})
	}

	out := make([]Out, len(result))
	for i, withReceptions := range result {
		out[i] = Out{
			PVZ:        withReceptions.PVZ,
			Receptions: withReceptions.Receptions,
		}
	}

	slog.DebugContext(ctx, "Usecase PVZ info success")
	return out, nil
}

func calcOffset(page, limit int) int {
	return (page - 1) * limit
}
