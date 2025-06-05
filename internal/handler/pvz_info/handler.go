package pvz_info

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/pvz_info"

	"github.com/go-playground/validator/v10"
)

type createHandler struct {
	usecase   usecase
	validator *validator.Validate
}

func New(usecase usecase, validator *validator.Validate) *createHandler {
	return &createHandler{
		usecase:   usecase,
		validator: validator,
	}
}

// @Summary Get PVZ information
// @Description Get PVZ list with receptions and products in given time period
// @ID GetPVZInfo
// @Tags PVZ
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body dto.PvzInfoIn true "Time period and pagination parameters"
// @Success 	200 {array} dto.PvzInfoOut "List of PVZs with receptions and products"
// @Failure 	400 {object} handler.errorResponse "Invalid parameters or no data found"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/pvz [get]
func (h *createHandler) GetPVZInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var request dto.PvzInfoIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if request.Page < 1 {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "incorrect page request", errors.New("page number must be greater than 0"))
		return
	}
	if request.Limit < 1 {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "incorrect limit request", errors.New("limit must be greater than 0"))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogStartDate(ctx, request.StartDate)
	ctx = logging.WithLogEndDate(ctx, request.EndDate)
	ctx = logging.WithLogPage(ctx, request.Page)
	ctx = logging.WithLogLimit(ctx, request.Limit)

	result, err := h.usecase.Run(ctx, pvz_info.In{
		StartData: request.StartDate,
		EndDate:   request.EndDate,
		Page:      request.Page,
		Limit:     request.Limit,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	out := make([]dto.PvzInfoOut, 0, len(result))
	for _, pvzWithReceptions := range result {
		pvz := dto.PvzInfoPvzOut{
			Uuid:             pvzWithReceptions.PVZ.Uuid,
			RegistrationDate: pvzWithReceptions.PVZ.RegistrationDate,
			City:             pvzWithReceptions.PVZ.City,
		}

		receptions := make([]dto.PvzInfoReceptionWithProductsOut, 0, len(pvzWithReceptions.Receptions))
		for _, recWithProds := range pvzWithReceptions.Receptions {
			reception := dto.PvzInfoReceptionOut{
				Id:       recWithProds.Reception.Uuid,
				DateTime: recWithProds.Reception.DateTime,
				PvzId:    recWithProds.Reception.PVZID,
				Status:   recWithProds.Reception.Status,
			}

			products := make([]dto.PvzInfoProductOut, 0, len(recWithProds.Products))
			for _, prod := range recWithProds.Products {
				products = append(products, dto.PvzInfoProductOut{
					Uuid:        prod.Uuid,
					DateTime:    prod.DateTime,
					Type:        prod.Type,
					ReceptionID: prod.ReceptionID,
				})
			}

			receptions = append(receptions, dto.PvzInfoReceptionWithProductsOut{
				Reception: reception,
				Products:  products,
			})
		}

		out = append(out, dto.PvzInfoOut{
			Pvz:        pvz,
			Receptions: receptions,
		})
	}
	if err = json.NewEncoder(w).Encode(out); err != nil {
		handler.RespondWithError(w, ctx, http.StatusInternalServerError, "failed to encode response", err)
		return
	}
}

func handleUseCaseError(w http.ResponseWriter, ctx context.Context, err error) {
	var statusCode int
	errorMsg := "internal server error"

	switch {
	case errors.Is(err, usecase2.ErrNotFoundProducts):
		errorMsg = "there is no products at this time interval"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrGetProducts):
		errorMsg = "failed to get products"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrGetReceptions):
		errorMsg = "failed to get receptions"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrGetPVZs):
		errorMsg = "failed to get pvz list"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrNotPageTooBig):
		errorMsg = "there is not enough pvzs to make such offset"
		statusCode = http.StatusBadRequest
	}

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
