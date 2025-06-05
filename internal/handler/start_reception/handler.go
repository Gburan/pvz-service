package start_reception

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/start_reception"

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

// @Summary Start new reception
// @Description Create new opened reception at PVZ
// @ID StartReception
// @Tags PVZ
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body dto.StartReceptionIn true "PVZ ID"
// @Success 	200 {object} dto.StartReceptionOut "Reception successfully started"
// @Failure 	400 {object} handler.errorResponse "PVZ not found or opened reception exists"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/receptions [post]
func (h *createHandler) StartReception(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var request dto.StartReceptionIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "failed to decode request", err)
		return
	}
	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogPVZID(ctx, request.PVZID)

	result, err := h.usecase.Run(ctx, start_reception.In{
		PVZID: request.PVZID,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	out := dto.StartReceptionOut{
		Id:       result.Reception.Uuid,
		DateTime: result.Reception.DateTime.UTC(),
		PvzId:    result.Reception.PVZID,
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
	case errors.Is(err, usecase2.ErrGetReception):
		errorMsg = "failed to get reception"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrFoundOpenedReception):
		errorMsg = "opened reception already exist"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrStartReception):
		errorMsg = "failed to start reception"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
