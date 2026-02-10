package server

import (
	"github.com/gin-gonic/gin"
	"rina.icu/pocman-server/internal/handler"
)

// Bind HTTP routes and return
func New(pingHandler *handler.PingHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", pingHandler.Ping)
	}

	return r
}
