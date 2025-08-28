// Package grpc used to receive comunication from other services
package grpc

import (
	"context"
	"errors"

	"github.com/FreyreCorona/Shortly/protos"
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/application"
	"github.com/FreyreCorona/Shortly/src/shortener_svc/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Service application.RetrieveURLService
	protos.UnimplementedGetURLServer
}

func NewGRPCServer(service application.RetrieveURLService) *Server {
	return &Server{Service: service}
}

func (s *Server) GetURLByShortCode(ctx context.Context, r *protos.GetURLRequest) (*protos.GetURLResponse, error) {
	url, err := s.Service.GetURL(r.GetCode())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCodeEmpty):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.NotFound, "the requested url doesnt exist")
		}
	}
	return &protos.GetURLResponse{RawUrl: url.RawURL, Code: url.ShortCode}, nil
}
