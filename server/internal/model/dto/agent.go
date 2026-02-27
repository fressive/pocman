package dto

import (
	"time"

	commonModel "github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/server/internal/data"
	"gorm.io/gorm"
)

// Agent represents a payload sent to CLI/HTTP consumers.
type Agent struct {
	gorm.Model   `json:"-"`
	AgentID      string `json:"agent_id" gorm:"size:64;not null;uniqueIndex"`
	CPUUsage     uint32 `json:"cpu_usage" gorm:"type:int;default:0"`
	RAMTotal     uint64 `json:"ram_total" gorm:"type:bigint;default:0"`
	RAMAvailable uint64 `json:"ram_available" gorm:"type:bigint;default:0"`

	LastInit      time.Time `json:"last_init"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

// Determine whether an agent is online or not.
func (a Agent) Online() bool {
	// Offline Threthold: 30s w/o heartbeat
	const OFFLINE_THRETHOLD = 30 * time.Second

	// if lastHB + 30s is before now, the agent goes offline
	return a.LastHeartbeat.Add(OFFLINE_THRETHOLD).Compare(time.Now()) != -1
}

func GetAgents() []Agent {
	var agents []Agent
	data.DB.Find(&agents)

	return agents
}

// ToModel converts the DTO back to the shared common Agent definition.
func (a Agent) ToModel() commonModel.Agent {
	return commonModel.Agent{
		AgentID:      a.AgentID,
		Online:       a.Online(),
		CPUUsage:     a.CPUUsage,
		RAMTotal:     a.RAMTotal,
		RAMAvailable: a.RAMAvailable,
		Uptime:       a.LastHeartbeat.Sub(a.LastInit).Seconds(),
	}
}
