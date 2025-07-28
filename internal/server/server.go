package server

import (
	"log"
	"net/http"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/logger"
	"oauth-tutorial/internal/mycrypto"
	authorize "oauth-tutorial/internal/presentation/authorization"
	"oauth-tutorial/internal/presentation/decision"
	tokenhandler "oauth-tutorial/internal/presentation/token"
	"oauth-tutorial/internal/session"
	usecase "oauth-tutorial/internal/usecase/authorization"
	u_decision "oauth-tutorial/internal/usecase/decision"
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
	ss := infrastructure.NewSessionStorage()
	acf := usecase.NewAuthorizationCodeFlow(s.logger, cr, sig, ss)

	rg := &mycrypto.RandomGenerator{}
	ur := infrastructure.NewUserRepository()
	ar := infrastructure.NewAuthCodeRepository()
	// 認可コード発行ユースケース
	pac := u_decision.NewPublishAuthorizationCodeUseCase(s.logger, rg, ss, ur, ar)

	// Set up HTTP handlers
	http.Handle("/authorize", authorize.NewAuthorizeHandler(s.logger, acf))
	http.Handle("/token", tokenhandler.NewTokenHandler(s.logger))
	http.Handle("/decision", decision.NewDecisionHandler(s.logger, pac))

	// Start the HTTP server
	s.logger.Info("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
