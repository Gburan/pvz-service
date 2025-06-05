package delete_product

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/delete_product"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

// @Summary Delete product
// @Description Delete product from PVZ. Requires JWT-Token with Employee role.
// @ID DeleteProduct
// @Tags PVZ
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param pvzId path string true "PVZ ID"
// @Success 	200 {object} object "Product successfully deleted"
// @Failure 	400 {object} handler.errorResponse "Invalid PVZ ID or no product to delete"
// @Failure 	401 {object} handler.errorResponse "Unauthorized"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/pvz/{pvzId}/delete_last_product [post]
func (h *createHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var request dto.DeleteProductIn
	vars := mux.Vars(r)
	pvzID, err := uuid.Parse(vars["pvzId"])
	if err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "invalid pvz format", err)
		return
	}
	request.PVZID = pvzID

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogPVZID(ctx, request.PVZID)

	err = h.usecase.Run(ctx, delete_product.In{
		PVZID: request.PVZID,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete product",
	}); err != nil {
		handler.RespondWithError(w, ctx, http.StatusInternalServerError, "failed to encode response", err)
		return
	}
}

func handleUseCaseError(w http.ResponseWriter, ctx context.Context, err error) {
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

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
