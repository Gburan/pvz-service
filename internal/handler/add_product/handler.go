package add_product

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
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

// @Summary Add a product to a PVZ reception
// @Description add product to PVZ's reception. Requires JWT-Token with Employee role.
// @ID AddProduct
// @Tags PVZ
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body dto.AddProductIn true "request params"
// @Success     200  {object}  dto.AddProductOut  "Product successfully added"
// @Failure     400  {object}  handler.errorResponse  "Invalid request body or validation failure"
// @Failure     401  {object}  handler.errorResponse  "Unauthorized - Invalid or missing JWT token"
// @Failure     500  {object}  handler.errorResponse  "Internal server error"
// @Router      /products [post]
func (h *createHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var request dto.AddProductIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogPVZID(ctx, request.PVZID)
	ctx = logging.WithLogProductType(ctx, request.Type)

	result, err := h.usecase.Run(ctx, add_product.In{
		Type:  request.Type,
		PVZID: request.PVZID,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	out := dto.AddProductOut{
		Uuid:        result.Product.Uuid,
		DateTime:    result.Product.DateTime.UTC(),
		Type:        result.Product.Type,
		ReceptionID: result.Product.ReceptionID,
	}
	if err = json.NewEncoder(w).Encode(out); err != nil {
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
	case errors.Is(err, usecase2.ErrAddProduct):
		errorMsg = "failed to add product"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
