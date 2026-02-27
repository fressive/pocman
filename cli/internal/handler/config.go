package handler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/fressive/pocman/cli/internal/conf"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

func promptConfig[T int | string | bool](name string, config *T) error {
	prompt := promptui.Prompt{
		Label:   name,
		Default: fmt.Sprint(*config),
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	switch pointer := any(config).(type) {
	case *int:
		value, parseErr := strconv.ParseInt(result, 0, 0)
		if parseErr != nil {
			return parseErr
		}
		*pointer = int(value)
	case *bool:
		value, parseErr := strconv.ParseBool(result)
		if parseErr != nil {
			return parseErr
		}
		*pointer = value
	case *string:
		*pointer = result
	}

	return nil
}

func ModifyConfig() error {
	// server section
	if err := promptConfig("[Server] Endpoint (e.g. http://127.0.0.1:5031)", &conf.CLIConfig.Server.Endpoint); err != nil {
		return err
	}
	if err := promptConfig("[Server] Token", &conf.CLIConfig.Server.Token); err != nil {
		return err
	}

	return nil
}

func InitConfig(ctx context.Context, c *cli.Command) error {
	fmt.Println("Initializing pocman-cli config...")

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
		if !c.Bool("override") {
			prompt := promptui.Prompt{
				Label:     "Config already exists, override",
				IsConfirm: true,
			}

			_, err := prompt.Run()

			if err != nil {
				return fmt.Errorf("Config already exists, quitted")
			}
		}

		// load this config and modify based on it
		prompt := promptui.Prompt{
			Label:     "Use the existing config as a template",
			IsConfirm: true,
		}

		_, err = prompt.Run()

		if err == nil {
			conf.CLIConfig.Load(configPath)
		}

	} else if !os.IsNotExist(err) {
		return err
	}

	ModifyConfig()

	prompt := promptui.Prompt{
		Label:     "Confirm",
		Default:   "Y",
		IsConfirm: true,
	}

	_, err = prompt.Run()

	if err == nil {
		conf.CLIConfig.Save(configPath)
	}
	return nil
}
