package http

import (
	"strings"

	"github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/model/dto"
	"github.com/fressive/pocman/server/internal/server/http/response"
	"github.com/gin-gonic/gin"
)

type VulnHandler struct{}

func NewVulnHandler() *VulnHandler {
	return &VulnHandler{}
}

func (h *VulnHandler) NewVuln(c *gin.Context) {
	var req model.Vuln
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 12001, "invalid request body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Code = strings.TrimSpace(req.Code)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		response.Error(c, 12002, "title is required")
		return
	}

	if req.Code == "" {
		response.Error(c, 12003, "code is required")
		return
	}

	record := dto.Vuln{
		Title:       req.Title,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := data.DB.Create(&record).Error; err != nil {
		if isUniqueConstraintErr(err) {
			response.Error(c, 12004, "vulnerability code already exists")
			return
		}
		response.Error(c, 12005, err.Error())
		return
	}

	response.Success(c, record.ToModel())
}

func isUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") ||
		strings.Contains(msg, "duplicate entry") ||
		strings.Contains(msg, "violates unique constraint")
}
