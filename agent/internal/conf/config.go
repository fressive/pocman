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
// prefix of "POCMANAGENT" and replaces dots in environment variable keys with
// underscores. If the configuration file cannot be read or unmarshalled,
// it returns an error.
package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Configuration of Pocman agent
type Config struct {
	// Agent name
	Name string `mapstructure:"name" yaml:"name"`

	// Mode of the agent, options: debug|release
	Mode string `mapstructure:"mode" yaml:"mode"`

	// Verbose mode, options: true|false
	Verbose bool `mapstructure:"verbose" yaml:"verbose"`

	Server *Server `mapstructure:"server" yaml:"server"`
	Data   *Data   `mapstructure:"data" yaml:"data"`
}

type Server struct {
	// Host of the server, example: 127.0.0.1
	Host *string `mapstructure:"host" yaml:"host"`

	// gRPC port of the server, example: 5032
	Port *int `mapstructure:"port" yaml:"port"`

	// Certificates to authenticate the identity
	GRPCCert *Certificate `mapstructure:"grpc_cert" yaml:"grpc_cert"`

	// Agent token (not recommended)
	GRPCToken *string `mapstructure:"grpc_token" yaml:"grpc_token"`
}

type Certificate struct {
	// Agent certificate file, exmaple: cert/agent.pem
	Cert string `mapstructure:"cert" yaml:"cert"`

	// Key file to certificate, example: cert/agent.key
	Key string `mapstructure:"key" yaml:"key"`

	// Self-signed CA certificate, example: cert/ca.pem
	CA string `mapstructure:"ca" yaml:"ca"`
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

var AgentConfig = Config{
	Mode:    "debug",
	Verbose: true,

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
	v.SetEnvPrefix("POCMANAGENT")
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
