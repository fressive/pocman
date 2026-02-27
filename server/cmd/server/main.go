package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/firebase/genkit/go/genkit"
	"github.com/fressive/pocman/server/internal/conf"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/llm"
	"github.com/fressive/pocman/server/internal/llm/tool"
	grpcServer "github.com/fressive/pocman/server/internal/server/grpc"
	httpServer "github.com/fressive/pocman/server/internal/server/http"
	"github.com/gin-gonic/gin"
)

var config_file = flag.String("c", "config.yml", "Configuration file, example: config.yml")

func main() {
	// parse cmd parameters
	flag.Parse()

	// init config
	slog.Info("Reading configuration file", "path", *config_file)
	if err := conf.ServerConfig.Load(*config_file); err != nil {
		slog.Error("failed to load config", "err", err)
		panic(err)
	}

	if conf.ServerConfig.Mode == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("You are running Pocman server with debug mode, set `mode` to `release` in your config to dismiss this warning")
	}

	// init Database
	err := data.InitDatabase()
	if err != nil {
		slog.Error("failed to init database", "err", err)
		panic(err)
	}

	// init LLM
	agent, err := llm.NewAgent()
	if err != nil {
		slog.Error("failed to init LLM Agent", "err", err)
		panic(err)
	}

	genkit.DefineTool(agent, "docker_compose_deploy", "Deploy containers from the docker-compose.yml file.", tool.DeployDockerCompose)

	// init HTTP server
	if conf.ServerConfig.Mode == "release" {
		// set gin log level
		gin.SetMode(gin.ReleaseMode)
	}

	httpSrv, err := httpServer.RunHTTPServer()
	if err != nil {
		slog.Error("failed to run HTTP server", "err", err)
		panic(err)
	}

	grpcSrv, err := grpcServer.RunGRPCServer()
	if err != nil {
		slog.Error("failed to run gRPC server", "err", err)
		panic(err)
	}

	// block and wait
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		slog.Error("Pocman server forced to shutdown.", "err", err)
	}

	slog.Info("Shutting down gRPC server...")
	grpcSrv.GracefulStop()

	slog.Info("Pocman server exited")
}
