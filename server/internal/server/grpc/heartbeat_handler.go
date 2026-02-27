package grpc

import (
	"context"
	"log/slog"
	"time"

	v1 "github.com/fressive/pocman/common/proto/v1"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/model/dto"
)

func (s *GRPCServer) Heartbeat(ctx context.Context, req *v1.HeartbeatRequest) (*v1.HeartbeatResponse, error) {
	slog.Debug("received heartbeat from agent", "ID", req.AgentId)

	data.DB.Where(&dto.Agent{
		AgentID: req.AgentId,
	}).Updates(dto.Agent{
		CPUUsage:      req.CpuUsage,
		RAMAvailable:  req.RamAvailable,
		LastHeartbeat: time.Now(),
	})

	return &v1.HeartbeatResponse{Status: "ok"}, nil
}
