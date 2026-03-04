package api

import (
	"context"

	"github.com/fressive/pocman/common/pkg/model"
)

func (c *Client) CreateVuln(ctx context.Context, vuln model.Vuln) (model.Vuln, error) {
	var created model.Vuln
	err := c.Do(ctx, "POST", "/api/v1/vuln", vuln, &created)
	if err != nil {
		return model.Vuln{}, err
	}

	return created, nil
}
