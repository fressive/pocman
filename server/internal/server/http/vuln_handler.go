package http

import (
	"errors"
	"strings"

	"github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/model/dto"
	"github.com/fressive/pocman/server/internal/server/http/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var existed dto.Vuln
	if err := data.DB.Where("code = ?", req.Code).First(&existed).Error; err == nil {
		response.Error(c, 12004, "vulnerability code already exists")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Error(c, 12005, err.Error())
		return
	}

	record := dto.Vuln{
		Title:       req.Title,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := data.DB.Create(&record).Error; err != nil {
		response.Error(c, 12005, err.Error())
		return
	}

	response.Success(c, record.ToModel())
}

func (h *VulnHandler) ListVuln(c *gin.Context) {

}
