package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/fressive/pocman/cli/internal"
	"github.com/fressive/pocman/cli/internal/conf"
	"github.com/fressive/pocman/cli/internal/handler"
	"github.com/urfave/cli/v3"
)

func readConfig(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// skip for "-v" flag
	if cmd.Bool("version") {
		return nil, nil
	}

	// Load CLI configuration
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
		conf.CLIConfig.Load(config)
	} else if os.IsNotExist(err) {
		fmt.Println("error: configuration file not found, use `pocman-cli config init` to initialize configuration")
		os.Exit(1)
		return nil, err
	} else {
		return nil, err
	}

	return nil, nil
}

func main() {
	cmd := &cli.Command{
		Name:    "pocman-cli",
		Version: internal.CLI_VERSION,
		Usage:   "a CLI tool for managing Pocman service",
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
					fmt.Printf("pocman-cli %s\n", internal.CLI_VERSION)
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
				Aliases: []string{"t"},
				Usage:   "Test the connection to the server",
				Before:  readConfig,
				Action:  handler.Test,
			},
			{
				Name:    "config",
				Aliases: []string{"conf"},
				Usage:   "Config operations",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "override",
						Aliases: []string{"r"},
						Usage:   "Force to override the existing configuration",
					},
				},
				Commands: []*cli.Command{
					{
						Name:    "init",
						Aliases: []string{"i"},
						Usage:   "Init the configuration",
						Action:  handler.InitConfig,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
