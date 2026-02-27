package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/fressive/pocman/common/proto/v1"
)

func (s *GRPCServer) Init(ctx context.Context, req *v1.InitRequest) (*v1.InitResponse, error) {
	slog.Info("new agent connected", "ID", req.AgentId, "version", req.Version)

	return &v1.InitResponse{Code: 0}, nil
}
