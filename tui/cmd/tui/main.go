package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/fressive/pocman/tui/internal"
	"github.com/fressive/pocman/tui/internal/command/agent"
	"github.com/fressive/pocman/tui/internal/command/config"
	"github.com/fressive/pocman/tui/internal/conf"
	"github.com/fressive/pocman/tui/internal/handler"
	"github.com/urfave/cli/v3"
)

func readConfig(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// skip for "-v" flag
	if cmd.Bool("version") {
		return nil, nil
	}

	// Load CLI config
	config := cmd.String("config")

	if config == "" {
		// use default path
		var err error
		config, err = conf.DefaultFilePath()

		if err != nil {
			return nil, err
		}
	}

	slog.Debug("using configuration", "file", config)

	if _, err := os.Stat(config); err == nil {
		conf.TUIConfig.Load(config)
	} else if os.IsNotExist(err) {
		// If config not exists, create it silently with default config
		conf.TUIConfig.Save(config)
	} else {
		return nil, err
	}

	return nil, nil
}

func main() {
	cmd := &cli.Command{
		Name:    "pocman-tui",
		Version: internal.CLI_VERSION,
		Usage:   "a CLI tool for managing Pocman service",
		Before:  readConfig,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Make the operation more talkative",
				Action: func(ctx context.Context, c *cli.Command, b bool) error {
					slog.SetLogLoggerLevel(slog.LevelDebug)
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"V"},
				Usage:   "Print version information and exit",
				Action: func(ctx context.Context, cmd *cli.Command, b bool) error {
					fmt.Printf("pocman-tui %s\n", internal.CLI_VERSION)
					os.Exit(0)
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "The location of config file, default to %UserConfigDir%/config.yml",
			},
		},

		Commands: []*cli.Command{
			{
				Name:    "test",
				Aliases: []string{"T"},
				Usage:   "Test the connection to the server",
				Action:  handler.Test,
			},
			{
				Name:    "agent",
				Aliases: []string{"a"},
				Usage:   "Agent related operations",
				Commands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "List all agents",
						Action:  agent.ListAgent,
					},
				},
			},
			{
				Name:    "config",
				Aliases: []string{"conf"},
				Usage:   "Configure pocman-tui tool",
				Action:  config.Configure,
			},
			{
				Name:    "vuln",
				Aliases: []string{"v"},
				Usage:   "Vulnerability related operations",
				Commands: []*cli.Command{
					{
						Name:    "new",
						Aliases: []string{"n"},
						Usage:   "Create a new vulnerability",
						Action:  handler.CreateVuln,
						Commands: []*cli.Command{
							{
								Name:   "cve",
								Usage:  "Create a new vulnerability from CVE",
								Action: handler.CreateVulnFromCVE,
							},
						},
					},
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "List remote vulnerabilities",
						Action:  handler.ListVuln,
					},
				},
			},
			{
				Name:    "task",
				Aliases: []string{"t"},
				Usage:   "Manage POC reproduction tasks",
				Commands: []*cli.Command{
					{
						Name:    "new",
						Aliases: []string{"n"},
						Usage:   "Create a new task",
						Action:  handler.NewTask,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println("pocman-tui: [ERROR] " + err.Error())
	}
}
