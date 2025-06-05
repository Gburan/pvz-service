package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"pvz-service/internal/logging"
)

type errorResponse struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func RespondWithError(w http.ResponseWriter, ctx context.Context, status int, errorMsg string, err error) {
	if status == http.StatusInternalServerError {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), fmt.Sprintf("Error: %s", err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := errorResponse{
		Message: errorMsg,
		Details: err.Error(),
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.ErrorContext(ctx, "Failed to encode error response", "error", err)
	}
}
