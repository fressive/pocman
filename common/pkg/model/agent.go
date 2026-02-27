package model

type Agent struct {
	AgentID      string `json:"agent_id"`
	Online       bool   `json:"online"`
	CPUUsage     uint32 `json:"cpu_usage"`
	RAMTotal     uint64 `json:"ram_total"`
	RAMAvailable uint64 `json:"ram_available"`
	Uptime       string `json:"uptime"`
}
