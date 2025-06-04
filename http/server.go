package http

import (
	"net/http"

	talenthub "github.com/caio86/talentHub"
)

type Server struct {
	server *http.Server
	router *http.ServeMux

	CandidatoService talenthub.CandidatoService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: http.NewServeMux(),
	}

	s.server.Handler = s.router

	newCandidatoHandler(
		s.CandidatoService,
	).loadRoutes(s.router)

	return s
}
