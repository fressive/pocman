package api

import (
	"context"
)

func (c *Client) Ping(ctx context.Context) error {
	return c.Do(ctx, "GET", "/api/v1/ping", nil, nil)
}
