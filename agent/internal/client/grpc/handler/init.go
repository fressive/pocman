package handler

import (
	"context"
	"log/slog"

	"github.com/fressive/pocman/agent/internal"
	"github.com/fressive/pocman/agent/internal/conf"
	v1 "github.com/fressive/pocman/common/proto/v1"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReportInit(client *v1.AgentServiceClient) error {
	hostInfo, err := host.Info()
	if err != nil {
		return err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	result, err := (*client).Init(context.Background(), &v1.InitRequest{
		Version:         internal.AGENT_VERSION,
		AgentId:         conf.AgentConfig.Name,
		Os:              hostInfo.OS,
		PlatformVersion: hostInfo.PlatformVersion,
		CpuModel:        cpuInfo[0].ModelName,
		CpuCores:        uint32(cpuInfo[0].Cores),
		CpuMhz:          uint32(cpuInfo[0].Mhz),
		RamTotal:        memInfo.Total,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				slog.Error("cannot authenticate the identity, check your credentials", "err", err)
			default:
				slog.Error("init report failed:", "err", err)
			}
		}
	}

	if result.Code != 0 {
		slog.Error("init failed", "code", result.Code)
		panic("failed to initialize agent")
	}

	slog.Info("report initialization successfully")
	return nil
}
