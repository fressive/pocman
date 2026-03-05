package api

import (
	"github.com/fressive/pocman/cli/internal"
	"github.com/fressive/pocman/cli/internal/conf"
	"github.com/fressive/pocman/common/pkg/api"
)

var apiClient *api.Client

func GetClient() (*api.Client, error) {
	if apiClient != nil {
		return apiClient, nil
	} else {
		var err error

		apiClient, err = api.NewClient(
			conf.CLIConfig.Server.Endpoint,
			conf.CLIConfig.Server.Token,
			internal.CLI_VERSION,
		)

		return apiClient, err
	}
}
