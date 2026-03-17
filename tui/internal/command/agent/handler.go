package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/fressive/pocman/tui/internal/api"
	"github.com/urfave/cli/v3"
)

func ListAgent(ctx context.Context, c *cli.Command) error {
	p := tea.NewProgram(initAgentModel(ctx))

	go func() {
		for {
			client, err := api.GetClient()

			if err != nil {
				p.Send(errMsg{err})
			}

			agents, err := client.ListAgents(ctx)

			if err != nil {
				p.Send(errMsg{err})
			}

			p.Send(agentMsg(agents))

			// based on the agent heartbeat report period (5s)
			time.Sleep(5 * time.Second)
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return nil
}
