package dummy_login

import (
	"encoding/json"
	"net/http"
	"time"

	"pvz-service/internal/handler"
	jwt2 "pvz-service/internal/jwt"

	"github.com/go-playground/validator/v10"
)

const (
	expIn   = 24 * time.Hour
	dummyId = "6c2f5cce-136d-4d80-9268-22abadc7bdf8"
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

func (h *createHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dummyLoginIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	token, err := jwt2.GenerateToken(h.secret, req.Role, dummyId, expIn)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "generate token failed", err)
		return
	}

	out := dummyLoginOut{
		Token: token,
	}
	if err = json.NewEncoder(w).Encode(out); err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "failed to encode response", err)
		return
	}
}
