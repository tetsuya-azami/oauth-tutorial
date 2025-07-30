package server

import (
	"log"
	"net/http"
	"oauth-tutorial/internal/infrastructure"
	pAuthorize "oauth-tutorial/internal/presentation/authorize"
	pDecision "oauth-tutorial/internal/presentation/decision"
	pToken "oauth-tutorial/internal/presentation/token"
	"oauth-tutorial/internal/session"
	uAuthorize "oauth-tutorial/internal/usecase/authorize"
	uDecision "oauth-tutorial/internal/usecase/decision"
	"oauth-tutorial/pkg/mycrypto"
	"oauth-tutorial/pkg/mylogger"
)

type Server struct {
	logger mylogger.Logger
}

func NewServer() *Server {
	return &Server{
		logger: mylogger.NewLogger(),
	}
}

func (s *Server) Start() {
	s.logger.Info("start server")

	// Initialize repositories and use cases
	cr := infrastructure.NewClientRepository()
	sig := session.NewSessionIDGenerator()
	ss := infrastructure.NewSessionStorage()
	acf := uAuthorize.NewAuthorizationCodeFlow(s.logger, cr, sig, ss)

	rg := &mycrypto.RandomGenerator{}
	ur := infrastructure.NewUserRepository()
	ar := infrastructure.NewAuthCodeRepository()
	// 認可コード発行ユースケース
	pac := uDecision.NewPublishAuthorizationCodeUseCase(s.logger, rg, ss, ur, ar)

	// Set up HTTP handlers
	http.Handle("/authorize", pAuthorize.NewAuthorizeHandler(s.logger, acf))
	http.Handle("/decision", pDecision.NewDecisionHandler(s.logger, pac))
	http.Handle("/token", pToken.NewTokenHandler(s.logger))

	// Start the HTTP server
	s.logger.Info("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
