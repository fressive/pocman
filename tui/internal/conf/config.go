package conf

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Server *Server `mapstructure:"server" yaml:"server"`
}

type Server struct {
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint"`

	Token string `mapstructure:"token" yaml:"token"`
}

var CLIConfig = Config{
	Server: &Server{
		Endpoint: "http://127.0.0.1:5031",
		Token:    "test-token",
	},
}

func DefaultDirPath() (string, error) {
	osConfigPath, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	configDir := path.Join(osConfigPath, "pocman_cli")

	return configDir, nil
}

func DefaultFilePath() (string, error) {
	configDirPath, err := DefaultDirPath()

	if err != nil {
		return "", err
	}

	configPath := path.Join(configDirPath, "config.yml")

	return configPath, nil
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
	v.SetEnvPrefix("POCMANCLI")
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

func (c *Config) Save(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	var m map[string]interface{}

	err := mapstructure.Decode(*c, &m)
	if err != nil {
		return err
	}

	err = v.MergeConfigMap(m)
	if err != nil {
		return err
	}

	return v.WriteConfigAs(path)
}
