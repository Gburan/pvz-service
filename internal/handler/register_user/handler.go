package register_user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"pvz-service/internal/handler"
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

func (h *registerHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request registerUserIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx := context.TODO()
	result, err := h.usecase.Run(ctx, register_user.In{
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	})
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	out := registerUserOut{
		Uuid:  result.User.ID,
		Email: result.User.Email,
		Role:  result.User.Role,
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

	handler.RespondWithError(w, statusCode, errorMsg, err)
}
