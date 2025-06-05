package close_reception

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/close_reception"

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

// @Summary Close opened reception at PVZ
// @Description Close current opened reception at PVZ. Requires JWT-Token with Employee role.
// @ID CloseReception
// @Tags PVZ
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param pvzId path string true "PVZ ID"
// @Success     200 {object} dto.CloseReceptionOut "Reception successfully closed"
// @Failure 	400 {object} handler.errorResponse "Invalid PVZ ID or no opened reception"
// @Failure 	401 {object} handler.errorResponse "Unauthorized"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/pvz/{pvzId}/close_last_reception [post]
func (h *createHandler) CloseReception(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var request dto.CloseReceptionIn
	vars := mux.Vars(r)
	pvzID, err := uuid.Parse(vars["pvzId"])
	if err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "invalid pvz format", err)
		return
	}
	request.PVZID = pvzID

	if err = h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogPVZID(ctx, request.PVZID)

	result, err := h.usecase.Run(ctx, close_reception.In{
		PVZID: request.PVZID,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	out := dto.CloseReceptionOut{
		Uuid:     result.Reception.Uuid,
		DateTime: result.Reception.DateTime.UTC(),
		PVZID:    result.Reception.PVZID,
		Status:   result.Reception.Status,
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
	case errors.Is(err, usecase2.ErrCloseReception):
		errorMsg = "failed to close reception"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
