// Package grpc iplements interface URLRepository
package grpc

import (
	"context"
	"time"

	"github.com/FreyreCorona/Shortly/protos"
	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCRepository struct {
	client protos.GetURLClient
}

func NewGRPCRepository(address string) (*GRPCRepository, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := protos.NewGetURLClient(conn)

	return &GRPCRepository{client: client}, nil
}

func (r *GRPCRepository) GetByShortCode(code string) (domain.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := protos.GetURLRequest{Code: code}

	res, err := r.client.GetURLByShortCode(ctx, &req)
	if err != nil {
		return domain.URL{}, err
	}

	return domain.URL{RawURL: res.RawUrl, ShortCode: res.Code}, nil
}
