package model

type Agent struct {
	AgentID      string  `json:"agent_id"`
	Online       bool    `json:"online"`
	CPUUsage     float32 `json:"cpu_usage"`
	RAMTotal     uint64  `json:"ram_total"`
	RAMAvailable uint64  `json:"ram_available"`
	Uptime       float64 `json:"uptime"`
	Load1        float32 `json:"load_1"`
	Load5        float32 `json:"load_5"`
	Load15       float32 `json:"load_15"`
}
