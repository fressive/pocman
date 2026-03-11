package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/fressive/pocman/agent/internal/conf"
	v1 "github.com/fressive/pocman/common/proto/v1"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReportHeartbeat(client *v1.AgentServiceClient) {
	for {
		memInfo, _ := mem.VirtualMemory()
		cpuUsage, _ := cpu.Percent(100*time.Millisecond, false)
		loadAvg, _ := load.Avg()

		_, err := (*client).Heartbeat(context.Background(), &v1.HeartbeatRequest{
			AgentId:      conf.AgentConfig.Name,
			CpuUsage:     float32(cpuUsage[0]),
			RamAvailable: memInfo.Available,
			Load_1:       float32(loadAvg.Load1),
			Load_5:       float32(loadAvg.Load5),
			Load_15:      float32(loadAvg.Load15),
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
