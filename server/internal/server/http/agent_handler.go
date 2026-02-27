package http

import (
	"github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/server/internal/model/dto"
	"github.com/fressive/pocman/server/internal/server/http/response"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AgentHandler struct{}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{}
}

func (h *AgentHandler) GetAgents(c *gin.Context) {
	agents := dto.GetAgents()
	response.Success(c, lo.Map(agents, func(a dto.Agent, _ int) model.Agent {
		return a.ToModel()
	}))
}
