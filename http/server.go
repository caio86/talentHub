package http

import "net/http"

type Server struct {
	server *http.Server
	router *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: http.NewServeMux(),
	}

	return s
}
