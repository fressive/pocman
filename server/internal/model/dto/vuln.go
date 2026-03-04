package dto

import (
	commonModel "github.com/fressive/pocman/common/pkg/model"
	"gorm.io/gorm"
)

// Vuln is the persisted vulnerability record exposed through HTTP APIs.
type Vuln struct {
	gorm.Model
	Title       string `json:"title" gorm:"size:256;not null"`
	Code        string `json:"code" gorm:"size:64;not null;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
}

func (v Vuln) ToModel() commonModel.Vuln {
	return commonModel.Vuln{
		ID:          uint32(v.ID),
		Title:       v.Title,
		Code:        v.Code,
		Description: v.Description,
	}
}
