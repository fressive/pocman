package grpc

import (
	"context"
	"log/slog"
	"time"

	v1 "github.com/fressive/pocman/common/proto/v1"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/model/dto"
)

func (s *GRPCServer) Init(ctx context.Context, req *v1.InitRequest) (*v1.InitResponse, error) {
	slog.Info("new agent connected", "ID", req.AgentId, "version", req.Version)

	var agent dto.Agent
	data.DB.Where(&dto.Agent{
		AgentID: req.AgentId,
	}).Assign(&dto.Agent{
		RAMTotal:      req.RamTotal,
		LastHeartbeat: time.Now(),
		LastInit:      time.Now(),
	}).FirstOrCreate(&agent)

	return &v1.InitResponse{Code: 0}, nil
}
