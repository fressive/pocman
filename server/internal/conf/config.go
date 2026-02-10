// Package conf provides functionality to load and manage application configuration.
//
// The Config struct holds the application configuration, including server,
// data, and logging settings. It is structured to be compatible with YAML
// configuration files and supports environment variable overrides.
//
// Config structure:
// - Server: Contains server-related settings such as Host, Port, and Mode.
// - Data: Contains database configuration settings.
// - Log: Contains logging configuration settings such as Level and Format.
//
// The Load function initializes the configuration by reading from a specified
// YAML file. It also supports automatic environment variable loading with a
// prefix of "POCMAN" and replaces dots in environment variable keys with
// underscores. If the configuration file cannot be read or unmarshalled,
// it returns an error.
package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Configuration of Pocman server
type Config struct {
	Server *Server `mapstructure:"server" yaml:"server"`
	Data   *Data   `mapstructure:"data" yaml:"data"`
}

type Server struct {
	// Host of the server, example: 127.0.0.1
	Host string `mapstructure:"host" yaml:"host"`

	// Port of the server, example: 5031
	Port int `mapstructure:"port" yaml:"port"`

	// Mode of the server, options: debug|release
	Mode string `mapstructure:"mode" yaml:"mode"`

	// Verbose mode, options: true|false
	Verbose bool `mapstructure:"verbose" yaml:"verbose"`
}

type Data struct {
	Database *Database `mapstructure:"database" yaml:"database"`
}

type Database struct {
	// Database driver, options: sqlite|mysql|postgres
	Driver string `mapstructure:"driver" yaml:"driver"`

	// DSN, example: file:data.db?cache=shared&mode=memory
	Source string `mapstructure:"source" yaml:"source"`
}

var AppConfig = Config{
	Server: &Server{
		Host:    "127.0.0.1",
		Port:    5031,
		Mode:    "debug",
		Verbose: true,
	},
	Data: &Data{
		Database: &Database{
			Driver: "sqlite",
			Source: "file:data.db?cache=shared&mode=memory",
		},
	},
}

func (c *Config) Load(path string) error {
	if path == "" {
		return fmt.Errorf("config path is empty")
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// find config in the environment
	v.AutomaticEnv()
	v.SetEnvPrefix("POCMAN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	v.WatchConfig()

	if err := v.Unmarshal(&c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
