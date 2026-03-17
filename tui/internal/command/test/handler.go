package test

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/urfave/cli/v3"
)

func TestConnection(ctx context.Context, c *cli.Command) error {
	p := tea.NewProgram(initTestModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return nil
}
