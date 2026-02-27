package api

import (
	"context"

	"github.com/fressive/pocman/common/pkg/model"
)

func (c *Client) Ping(ctx context.Context) error {
	var resp model.Response[any]
	err := c.Do(ctx, "GET", "/api/v1/ping", nil, resp)
	return err
}
