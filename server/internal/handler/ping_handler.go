package handler

import (
	"github.com/fressive/pocman/server/internal/server/response"
	"github.com/gin-gonic/gin"
)

// PingHandler handles ping endpoints for health checks.
type PingHandler struct{}

// NewPingHandler creates a new PingHandler.
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping responds with a simple pong payload.
func (h *PingHandler) Ping(c *gin.Context) {
	response.Success(c, gin.H{"ping": "pong"})
}
