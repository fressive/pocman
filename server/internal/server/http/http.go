package http

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/fressive/pocman/server/internal/conf"
	"github.com/fressive/pocman/server/internal/model/dto"
	"github.com/fressive/pocman/server/internal/server/http/response"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")

		if auth != "" {
			// Bearer xxxxxxxx
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
				response.Unauth(ctx, nil)
				ctx.Abort()
				return
			}

			b64token := parts[1]
			tokenb, err := base64.RawURLEncoding.DecodeString(b64token)

			if err != nil {
				response.Error(ctx, 20001, "fail to decode token")
				ctx.Abort()
				return
			}

			token := string(tokenb)

			if err := dto.VerifyToken(token); err == nil {
				ctx.Next()
				return
			} else {
				response.Unauth(ctx, err)
				ctx.Abort()
				return
			}
		}

		response.Unauth(ctx, nil)
		ctx.Abort()
	}
}

// Bind HTTP routes and return
func NewHTTPServer(pingHandler *PingHandler, agentHandler *AgentHandler, fileHandler *FileHandler, vulnHandler *VulnHandler) (*gin.Engine, error) {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(TokenAuthMiddleware())

	if conf.ServerConfig.Mode == "debug" {
		r.Use(gin.Logger())
	}

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
