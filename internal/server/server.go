package server

import (
	"log"
	"net/http"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/logger"
	authorize "oauth-tutorial/internal/presentation/authorization"
	tokenhandler "oauth-tutorial/internal/presentation/token"
	"oauth-tutorial/internal/session"
	usecase "oauth-tutorial/internal/usecase/authorization"
)

type Server struct {
	logger logger.MyLogger
}

func NewServer() *Server {
	return &Server{
		logger: logger.NewMyLogger(),
	}
}

func (s *Server) Start() {
	s.logger.Info("start server")

	// Initialize repositories and use cases
	cr := infrastructure.NewClientRepository()
	sig := session.NewSessionIDGenerator()
	aps := infrastructure.NewAuthParamSession()
	acf := usecase.NewAuthorizationCodeFlow(s.logger, cr, sig, aps)

	// Set up HTTP handlers
	http.Handle("/authorize", authorize.NewAuthorizeHandler(s.logger, acf))
	http.Handle("/token", tokenhandler.NewTokenHandler(s.logger))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
