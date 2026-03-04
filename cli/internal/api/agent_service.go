package api

import (
	"context"

	"github.com/fressive/pocman/common/pkg/model"
)

func (c *Client) ListAgents(ctx context.Context) ([]model.Agent, error) {
	var agents []model.Agent
	err := c.Do(ctx, "GET", "/api/v1/agent", nil, &agents)

	if err != nil {
		return nil, err
	}

	return agents, err
}
