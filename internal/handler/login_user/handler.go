package login_user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	jwt2 "pvz-service/internal/jwt"
	"pvz-service/internal/logging"
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

// @Summary User login
// @Description Authenticate user and get JWT token
// @ID LoginUser
// @Tags User
// @Accept json
// @Produce json
// @Param input body dto.LoginUserIn true "Login credentials"
// @Success 	200 {object} dto.LoginUserOut "Successfully logged in"
// @Failure 	400 {object} handler.errorResponse "Invalid credentials or validation failed"
// @Failure 	401 {object} handler.errorResponse "Unauthorized"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/login [post]
func (h *loginHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var request dto.LoginUserIn
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	if len(request.Email) > 4 {
		email := strings.Repeat("*", len(request.Email)-4) + request.Email[len(request.Email)-4:]
		ctx = logging.WithLogEmail(ctx, email)
	}

	result, err := h.usecase.Run(ctx, login_user.In{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		handleUseCaseError(w, ctx, err)
		return
	}

	token, err := jwt2.GenerateToken(h.secret, result.User.Role, result.User.Uuid, expIn)
	if err != nil {
		handler.RespondWithError(w, ctx, http.StatusInternalServerError, "generate token failed", err)
		return
	}

	out := dto.LoginUserOut{
		Token: token,
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

	handler.RespondWithError(w, ctx, statusCode, errorMsg, err)
}
