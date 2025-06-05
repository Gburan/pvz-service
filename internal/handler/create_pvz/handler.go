package create_pvz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/create_pvz"

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

// @Summary Create PVZ
// @Description Create new PVZ. Requires JWT-Token with Employee role.
// @ID CreatePVZ
// @Tags PVZ
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body dto.CreatePVZIn true "PVZ creation data"
// @Success 	200 {object} dto.CreatePVZOut "PVZ successfully created"
// @Failure 	400 {object} handler.errorResponse "Invalid request or validation failed"
// @Failure 	401 {object} handler.errorResponse "Unauthorized"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/pvz [post]
func (h *createHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var request dto.CreatePVZIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogCity(ctx, request.City)

	result, err := h.usecase.Run(ctx, create_pvz.In{
		City: request.City,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	out := dto.CreatePVZOut{
		Uuid:             result.PVZ.Uuid,
		RegistrationDate: result.PVZ.RegistrationDate.UTC(),
		City:             result.PVZ.City,
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
	case errors.Is(err, usecase2.ErrAddPVZ):
		errorMsg = "failed to add pvz"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
