package delete_product

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"pvz-service/internal/handler"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/delete_product"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

func (h *createHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var request deleteProductIn
	vars := mux.Vars(r)
	request.PVZID = vars["pvzId"]

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx := context.TODO()
	err := h.usecase.Run(ctx, delete_product.In{
		PVZID: request.PVZID,
	})
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete product",
	}); err != nil {
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
	case errors.Is(err, usecase2.ErrNotFoundProduct):
		errorMsg = "no product to delete"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrGetProduct):
		errorMsg = "failed to find product to delete"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrDeleteProduct):
		errorMsg = "failed to delete product"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, statusCode, errorMsg, err)
}
