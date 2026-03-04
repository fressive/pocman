package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fressive/pocman/server/internal/conf"
	"github.com/gin-gonic/gin"
)

// Bind HTTP routes and return
func NewHTTPServer(pingHandler *PingHandler, agentHandler *AgentHandler, fileHandler *FileHandler, vulnHandler *VulnHandler) (*gin.Engine, error) {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", pingHandler.Ping)

		v1.GET("/agent", agentHandler.GetAgents)

		v1.POST("/vuln", vulnHandler.NewVuln)

		v1.POST("/file/upload", fileHandler.FileUpload)
		v1.GET("/file/download", fileHandler.FileDownload)
	}

	return r, nil
}

func RunHTTPServer() (*http.Server, error) {
	pingHandler := NewPingHandler()
	agentHandler := NewAgentHandler()
	fileHandler := NewFileHandler()
	vulnHandler := NewVulnHandler()

	r, err := NewHTTPServer(pingHandler, agentHandler, fileHandler, vulnHandler)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", conf.ServerConfig.Server.Host, conf.ServerConfig.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		slog.Info("Starting Pocman HTTP server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error occured when HTTP serves", "err", err)
		}
	}()

	return srv, nil
}
