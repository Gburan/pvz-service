package login_user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"pvz-service/internal/handler"
	jwt2 "pvz-service/internal/jwt"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/login_user"

	"github.com/go-playground/validator/v10"
)

const (
	expIn = 24 * time.Hour
)

type loginHandler struct {
	secret    string
	usecase   usecase
	validator *validator.Validate
}

func New(secret string, usecase usecase, validator *validator.Validate) *loginHandler {
	return &loginHandler{
		secret:    secret,
		usecase:   usecase,
		validator: validator,
	}
}

func (h *loginHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request loginUserIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	ctx := context.TODO()
	result, err := h.usecase.Run(ctx, login_user.In{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	token, err := jwt2.GenerateToken(h.secret, result.User.Role, result.User.ID, expIn)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "generate token failed", err)
		return
	}

	out := loginUserOut{
		Token: token,
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
	case errors.Is(err, usecase2.ErrNotFoundUser):
		errorMsg = "not found such user"
		statusCode = http.StatusBadRequest
	case errors.Is(err, usecase2.ErrGetUser):
		errorMsg = "error while get user data"
		statusCode = http.StatusInternalServerError
	case errors.Is(err, usecase2.ErrIncorrectPass):
		errorMsg = "incorrect password"
		statusCode = http.StatusBadRequest
	}

	handler.RespondWithError(w, statusCode, errorMsg, err)
}
