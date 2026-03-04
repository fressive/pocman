package handler

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/fressive/pocman/cli/internal/conf"
	"github.com/urfave/cli/v3"
)

func generateField(name string, config any) huh.Field {
	switch config := config.(type) {
	case *string:
		return huh.NewInput().Title(name).Value(config)
	case *bool:
		return huh.NewConfirm().Title(name).Value(config).Affirmative("true").Negative("false")
	}

	return nil
}

func ModifyConfig() error {
	// server section

	form := huh.NewForm(
		huh.NewGroup(
			generateField("Endpoint", &conf.CLIConfig.Server.Endpoint),
			generateField("Token", &conf.CLIConfig.Server.Token),
		).Title("Server"),
	)

	return form.Run()
}

func Configure(ctx context.Context, c *cli.Command) error {
	configDir, err := conf.DefaultDirPath()

	if err != nil {
		return err
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		slog.Debug("Default config path not exist, creating dir...", "path", configDir)
		err = os.Mkdir(configDir, os.ModePerm)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	configPath, err := conf.DefaultFilePath()

	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); err == nil {
		conf.CLIConfig.Load(configPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	err = ModifyConfig()
	if err == nil {
		conf.CLIConfig.Save(configPath)
	}

	return err
}
