package handler

import (
	"context"
	"fmt"

	"github.com/fressive/pocman/tui/internal/api"
	"github.com/urfave/cli/v3"
)

func Test(ctx context.Context, c *cli.Command) error {
	fmt.Println("pocman-tui: Testing connection...")

	client, err := api.GetClient()
	if err != nil {
		return err
	}

	err = client.Ping(ctx)
	if err != nil {
		fmt.Println("pocman-tui: Failed to connect. Check the configuration by using `pocman-tui config`")
		return err
	}

	fmt.Println("pocman-tui: Connect successfully.")
	return nil
}
