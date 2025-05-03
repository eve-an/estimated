package api

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
	handler    http.Handler
}

func NewServer(
	httpServer *http.Server,
	handler http.Handler,
) *Server {
	return &Server{
		httpServer: httpServer,
		handler:    handler,
	}
}

func (s *Server) Serve() error {
	s.httpServer.Handler = s.handler
	return s.httpServer.ListenAndServe()
}
