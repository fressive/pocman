package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/fressive/pocman/common/proto/v1"
)

func (s *GRPCServer) Heartbeat(ctx context.Context, req *v1.HeartbeatRequest) (*v1.HeartbeatResponse, error) {
	slog.Debug("received heartbeat from agent", "ID", req.AgentId)

	return &v1.HeartbeatResponse{Status: "ok"}, nil
}
