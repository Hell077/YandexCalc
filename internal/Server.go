package internal

import (
	"net/http"
)

type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/api/v1/calculate", CalculateHandler)
	return http.ListenAndServe(s.address, nil)
}
