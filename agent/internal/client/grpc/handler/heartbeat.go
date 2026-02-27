package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/fressive/pocman/agent/internal/conf"
	v1 "github.com/fressive/pocman/common/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReportHeartbeat(client *v1.AgentServiceClient) {
	for {
		_, err := (*client).Heartbeat(context.Background(), &v1.HeartbeatRequest{
			AgentId:  conf.AgentConfig.Name,
			CpuUsage: 12.5,
		})

		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.Unauthenticated:
					slog.Error("cannot authenticate the identity, check your credentials", "err", err)
				default:
					slog.Error("heartbeat report failed:", "err", err)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
