package http

import (
	"net/http"

	talenthub "github.com/caio86/talentHub"
)

type Server struct {
	server *http.Server
	router *http.ServeMux

	Addr string

	// Services
	CandidatoService talenthub.CandidatoService
	VagaService      talenthub.VagaService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: http.NewServeMux(),
	}

	router := http.NewServeMux()

	// Loading routes
	s.loadCandidatoRoutes(router)
	s.loadVagaRoutes(router)

	s.router.Handle("/api/v1/", http.StripPrefix("/api/v1", router))
	// Setting middlewares
	middlewares := createMiddlewares(
		s.logging,
	)

	s.server.Handler = middlewares(s.router)

	return s
}

func (s *Server) Open() error {
	s.server.Addr = s.Addr

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
