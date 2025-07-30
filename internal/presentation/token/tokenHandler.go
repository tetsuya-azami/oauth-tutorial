package tokenhandler

import (
	"net/http"
	"oauth-tutorial/pkg/mylogger"
)

type TokenHandler struct {
	logger mylogger.Logger
}

func NewTokenHandler(logger mylogger.Logger) *TokenHandler {
	return &TokenHandler{logger: logger}
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("token endpoint called")
	w.Write([]byte("token endpoint"))
}
