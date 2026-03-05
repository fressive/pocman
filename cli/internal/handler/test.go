package handler

import (
	"context"
	"fmt"

	"github.com/fressive/pocman/cli/internal/api"
	"github.com/urfave/cli/v3"
)

func Test(ctx context.Context, c *cli.Command) error {
	fmt.Println("pocman-cli: Testing connection...")

	client, err := api.GetClient()
	if err != nil {
		return err
	}

	err = client.Ping(ctx)
	if err != nil {
		fmt.Println("pocman-cli: Failed to connect. Check the configuration by using `pocman-cli config`")
		return err
	}

	fmt.Println("pocman-cli: Connect successfully.")
	return nil
}
