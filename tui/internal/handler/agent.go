package handler

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/tui/internal/api"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
)

func ListAgents(ctx context.Context, c *cli.Command) error {
	client, err := api.GetClient()
	if err != nil {
		return err
	}

	agents, err := client.ListAgents(ctx)

	if err != nil {
		return err
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint()),
		tablewriter.WithRendition(tw.Rendition{
			Borders: tw.BorderNone,
			Symbols: tw.NewSymbols(tw.StyleNone),
			Settings: tw.Settings{
				Lines:      tw.LinesNone,
				Separators: tw.SeparatorsNone,
			},
		}),
	)
	table.Header("ID", "Status", "Uptime", "CPU", "RAM")
	table.Bulk(lo.Map(agents, func(a model.Agent, _ int) []any {
		var online string
		var uptime string

		if a.Online {
			online = "Online"
			uptimeDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", math.Round(a.Uptime)))
			uptime = uptimeDuration.String()
		} else {
			online = "Offline"
			uptime = "N/A"
		}

		return []any{
			a.AgentID,
			online,
			uptime,
			fmt.Sprintf("%.0f%%", a.CPUUsage),
			fmt.Sprintf("%dM/%dM", (a.RAMTotal-a.RAMAvailable)/1024/1024, a.RAMTotal/1024/1024),
		}
	}))
	table.Render()

	return nil
}
