package add_product

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"pvz-service/internal/handler"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/add_product"

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

func (h *createHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request addProductIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx := context.TODO()
	result, err := h.usecase.Run(ctx, add_product.In{
		Type:  request.Type,
		PVZID: request.PVZID,
	})
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	out := addProductOut{
		Uuid:        result.Product.Uuid,
		DateTime:    result.Product.DateTime.UTC().Format(time.RFC3339Nano),
		Type:        result.Product.Type,
		ReceptionID: result.Product.ReceptionID,
	}
	if err = json.NewEncoder(w).Encode(out); err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "failed to encode response", err)
		return
	}
}

func handleUseCaseError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	errorMsg := "internal server error"

	switch {
	case errors.Is(err, usecase2.ErrNotFoundPVZ):
		errorMsg = "pvz with such id not exist"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrGetPVZByID):
		errorMsg = "failed to look up for pvz"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrNotFoundReception):
		errorMsg = "there is no receptions at all"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrGetReception):
		errorMsg = "failed to get reception"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrNotFoundOpenedReception):
		errorMsg = "no opened reception"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrAddProduct):
		errorMsg = "failed to add product"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, statusCode, errorMsg, err)
}
