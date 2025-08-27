// Package grpc used to receive comunication from other services
package grpc

import (
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/application"
)

type Server struct {
	Service application.RetrieveURLService
}

func NewGRPCServer(service application.RetrieveURLService) *Server {
	return &Server{Service: service}
}

func (s *Server) RetrieveURL(code string) error {
	return nil
}
