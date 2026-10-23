package main

import "net/http"

type Server struct {
	kvs *KVS
	mux *http.ServeMux
}

func NewServer(kvs *KVS) *Server {
	s := &Server{
		kvs: kvs,
		mux: http.NewServeMux(),
	}

	s.registeredRoute()

	return s
}

func (s *Server) registeredRoute() {
	s.mux.HandleFunc("GET /health", s.healthHandler)

	s.mux.HandleFunc("GET /cache/{key}", s.getHandler)
	s.mux.HandleFunc("PUT /cache/{key}", s.setHandler)
	s.mux.HandleFunc("DELETE /cache/{key}", s.delHandler)
	s.mux.HandleFunc("GET /pong", s.pingPong)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
