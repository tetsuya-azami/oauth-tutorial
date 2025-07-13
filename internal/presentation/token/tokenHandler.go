package tokenhandler

import (
	"log/slog"
	"net/http"
)

type TokenHandler struct {
	logger *slog.Logger
}

func NewTokenHandler(logger *slog.Logger) *TokenHandler {
	return &TokenHandler{logger: logger}
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("token endpoint called")
	w.Write([]byte("token endpoint"))
}
