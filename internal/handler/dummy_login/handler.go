package dummy_login

import (
	"encoding/json"
	"net/http"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	"pvz-service/internal/handler"
	jwt2 "pvz-service/internal/jwt"
	"pvz-service/internal/logging"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	expIn = 24 * time.Hour
)

var (
	dummyId = uuid.New()
)

type createHandler struct {
	secret    string
	validator *validator.Validate
}

func New(secret string, validator *validator.Validate) *createHandler {
	return &createHandler{
		secret:    secret,
		validator: validator,
	}
}

// @Summary Dummy login
// @Description Get JWT token for testing purposes
// @ID DummyLogin
// @Tags User
// @Accept json
// @Produce json
// @Param input body dto.DummyLoginIn true "Login credentials"
// @Success 	200 {object} dto.DummyLoginOut "Successfully logged in"
// @Failure 	400 {object} handler.errorResponse "Invalid request or validation failed"
// @Failure 	500 {object} handler.errorResponse "Internal server error"
// @Router 		/dummyLogin [post]
func (h *createHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req dto.DummyLoginIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		handler.RespondWithError(w, ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	token, err := jwt2.GenerateToken(h.secret, req.Role, dummyId, expIn)
	if err != nil {
		handler.RespondWithError(w, ctx, http.StatusInternalServerError, "generate token failed", err)
		return
	}

	ctx = logging.WithLogRole(ctx, req.Role)
	out := dto.DummyLoginOut{
		Token: token,
	}
	if err = json.NewEncoder(w).Encode(out); err != nil {
		handler.RespondWithError(w, ctx, http.StatusInternalServerError, "failed to encode response", err)
		return
	}
}
