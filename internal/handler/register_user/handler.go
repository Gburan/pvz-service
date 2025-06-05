package register_user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	"pvz-service/internal/logging"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/register_user"

	"github.com/go-playground/validator/v10"
)

type registerHandler struct {
	usecase   usecase
	validator *validator.Validate
}

func New(usecase usecase, validator *validator.Validate) *registerHandler {
	return &registerHandler{
		usecase:   usecase,
		validator: validator,
	}
}

// @Summary Register new user
// @Description Create new user account
// @ID RegisterUser
// @Tags User
// @Accept json
// @Produce json
// @Param input body dto.RegisterUserIn true "User registration data"
// @Success 	200 {object} dto.RegisterUserOut "User successfully registered"
// @Failure 	400 {object} handler.errorResponse "Validation failed or user already exists"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/register [post]
func (h *registerHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var request dto.RegisterUserIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx = logging.WithLogEmail(ctx, request.Email)
	ctx = logging.WithLogRole(ctx, request.Role)

	result, err := h.usecase.Run(ctx, register_user.In{
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	out := dto.RegisterUserOut{
		Uuid:  result.User.Uuid,
		Email: result.User.Email,
		Role:  result.User.Role,
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
	case errors.Is(err, usecase2.ErrGetUser):
		errorMsg = "failed to look up for existing user"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrUserAlreadyExist):
		errorMsg = "such user already exist"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrGenHashedPass):
		errorMsg = "error while gen has password"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrAddUser):
		errorMsg = "error while adding user in db"
		statusCode = http.StatusInternalServerError
	}

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
