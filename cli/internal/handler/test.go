package handler

import (
	"context"
	"fmt"

	"github.com/fressive/pocman/cli/internal/api"
	"github.com/urfave/cli/v3"
)

func Test(ctx context.Context, c *cli.Command) error {
	fmt.Println("Testing connection...")

	client, err := api.GetClient()
	if err != nil {
		return err
	}

	err = client.Ping(ctx)
	if err != nil {
		fmt.Println("Failed to connect. Check the configuration by using `pocman-cli config list`")
		return err
	}

	fmt.Println("Connect successfully.")
	return nil
}
