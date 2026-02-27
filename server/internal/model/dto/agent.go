package dto

import (
	commonModel "github.com/fressive/pocman/common/pkg/model"
	"gorm.io/gorm"
)

// Agent represents a payload sent to CLI/HTTP consumers.
type Agent struct {
	gorm.Model   `json:"-"`
	AgentID      string `json:"agent_id" gorm:"size:64;not null;uniqueIndex"`
	Online       bool   `json:"online" gorm:"default:false"`
	CPUUsage     uint32 `json:"cpu_usage" gorm:"type:int;default:0"`
	RAMTotal     uint64 `json:"ram_total" gorm:"type:bigint;default:0"`
	RAMAvailable uint64 `json:"ram_available" gorm:"type:bigint;default:0"`
	Uptime       string `json:"uptime" gorm:"size:64"`
}

// FromModel converts the shared common Agent into DTO form.
func FromModel(src commonModel.Agent) Agent {
	return Agent{
		AgentID:      src.AgentID,
		Online:       src.Online,
		CPUUsage:     src.CPUUsage,
		RAMTotal:     src.RAMTotal,
		RAMAvailable: src.RAMAvailable,
		Uptime:       src.Uptime,
	}
}

// ToModel converts the DTO back to the shared common Agent definition.
func (a Agent) ToModel() commonModel.Agent {
	return commonModel.Agent{
		AgentID:      a.AgentID,
		Online:       a.Online,
		CPUUsage:     a.CPUUsage,
		RAMTotal:     a.RAMTotal,
		RAMAvailable: a.RAMAvailable,
		Uptime:       a.Uptime,
	}
}
