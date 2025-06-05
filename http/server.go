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

	// Loading routes
	s.loadCandidatoRoutes(s.router)
	s.loadVagaRoutes(s.router)

	middlewares := createMiddlewares(
		logging,
	)

	router := middlewares(s.router)

	s.server.Handler = router

	return s
}

func (s *Server) Open() error {
	s.server.Addr = s.Addr

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
