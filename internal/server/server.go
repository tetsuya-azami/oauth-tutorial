package server

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/logger"
)

type Hoge struct {
	Id   string
	Name string
	Arr  []string
}

type Server struct {
	logger *slog.Logger
}

func NewServer() *Server {
	return &Server{
		logger: logger.NewMyLogger(),
	}
}

func (s *Server) Start() {
	s.logger.Info("start server")
	http.HandleFunc("GET /authorize", s.authorizeHandler)
	http.HandleFunc("POST /token", s.tokenHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *Server) authorizeHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.validateAuthorizeRequest(r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := "dummy-auth-code"

	s.logger.Info("Authorization code generated", "code", code)
	// http.Redirect(w, r, req.RedirectURI()+"?code="+code+"&state="+req.State(), http.StatusFound)
}

func (s *Server) validateAuthorizeRequest(queries url.Values) error {
	// 認可フローによって処理が異なるケースを想定。現状、認可コードフローのみサポートのため、取得した値を使用してはいない
	_, err := domain.GetResponseType(queries.Get("response_type"))
	if err != nil {
		s.logger.Info("Invalid authorize request parameters", "error", err)
		return err
	}

	if queries.Get("client_id") == "" {
		s.logger.Info("client_id is empty")
		return errors.New("client_id is required")
	}

	return nil
}

func (s *Server) tokenHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("token endpoint called")
	w.Write([]byte("token endpoint"))
}
