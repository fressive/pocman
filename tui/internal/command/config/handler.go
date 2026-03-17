package config

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/urfave/cli/v3"
)

func Configure(ctx context.Context, c *cli.Command) error {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	return nil
}
