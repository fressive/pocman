package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	client "github.com/fressive/pocman/agent/internal/client/grpc"
	"github.com/fressive/pocman/agent/internal/client/grpc/handler"
	"github.com/fressive/pocman/agent/internal/conf"
	protocol "github.com/fressive/pocman/common/proto/v1"
)

var configFile = flag.String("c", "config.yml", "Configuration file, example: config.yml")

func main() {
	// parse cmd parameters
	flag.Parse()

	// reading config
	slog.Info("Reading configuration file", "path", *configFile)
	if err := conf.AgentConfig.Load(*configFile); err != nil {
		slog.Error("failed to load config", "err", err)
		panic(err)
	}

	// create gRPC client
	conn, err := client.NewConn()
	if err != nil {
		slog.Error("failed to create gRPC client", "err", err)
		panic(err)
	}

	conn.Connect()
	defer conn.Close()

	c := protocol.NewAgentServiceClient(conn)
	handler.ReportInit(&c)

	go handler.ReportHeartbeat(&c)

	// block and wait
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Pocman agent exited")
}
