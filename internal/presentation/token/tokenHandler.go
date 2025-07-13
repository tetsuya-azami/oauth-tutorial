package tokenhandler

import (
	"net/http"
	"oauth-tutorial/internal/logger"
)

type TokenHandler struct {
	logger logger.MyLogger
}

func NewTokenHandler(logger logger.MyLogger) *TokenHandler {
	return &TokenHandler{logger: logger}
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("token endpoint called")
	w.Write([]byte("token endpoint"))
}
