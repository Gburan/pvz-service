package pvz_info

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"pvz-service/internal/handler"
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

func (h *createHandler) GetPVZInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request pvzInfoIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if request.Page < 1 {
		handler.RespondWithError(w, http.StatusBadRequest, "incorrect page request", nil)
		return
	}
	if request.Limit < 1 {
		handler.RespondWithError(w, http.StatusBadRequest, "incorrect limit request", nil)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx := context.TODO()
	result, err := h.usecase.Run(ctx, pvz_info.In{
		StartData: request.StartDate,
		EndDate:   request.EndDate,
		Page:      request.Page,
		Limit:     request.Limit,
	})
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	out := make([]pvzInfoOut, 0, len(result))
	for _, pvzWithReceptions := range result {
		pvz := pvzOut{
			Uuid:             pvzWithReceptions.PVZ.Uuid,
			RegistrationDate: pvzWithReceptions.PVZ.RegistrationDate.Format(time.RFC3339),
			City:             pvzWithReceptions.PVZ.City,
		}

		receptions := make([]receptionWithProductsOut, 0, len(pvzWithReceptions.Receptions))
		for _, recWithProds := range pvzWithReceptions.Receptions {
			reception := receptionOut{
				Id:       recWithProds.Reception.Uuid,
				DateTime: recWithProds.Reception.DateTime.Format(time.RFC3339),
				PvzId:    recWithProds.Reception.PVZID,
				Status:   recWithProds.Reception.Status,
			}

			products := make([]productOut, 0, len(recWithProds.Products))
			for _, prod := range recWithProds.Products {
				products = append(products, productOut{
					Uuid:        prod.Uuid,
					DateTime:    prod.DateTime.Format(time.RFC3339),
					Type:        prod.Type,
					ReceptionID: prod.ReceptionID,
				})
			}

			receptions = append(receptions, receptionWithProductsOut{
				Reception: reception,
				Products:  products,
			})
		}

		out = append(out, pvzInfoOut{
			Pvz:        pvz,
			Receptions: receptions,
		})
	}
	if err = json.NewEncoder(w).Encode(out); err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "failed to encode response", err)
		return
	}
}

func handleUseCaseError(w http.ResponseWriter, err error) {
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

	handler.RespondWithError(w, statusCode, errorMsg, err)
}
