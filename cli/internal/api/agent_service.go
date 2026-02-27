package api

import (
	"context"

	"github.com/fressive/pocman/common/pkg/model"
)

func (c *Client) ListAgents(ctx context.Context) ([]model.Agent, error) {
	var resp model.Response[[]model.Agent]
	err := c.Do(ctx, "GET", "/api/v1/agent", nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Data, err
}
