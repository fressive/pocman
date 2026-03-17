package api

import (
	"github.com/fressive/pocman/common/pkg/api"
	"github.com/fressive/pocman/tui/internal"
	"github.com/fressive/pocman/tui/internal/conf"
)

var apiClient *api.Client

func GetClient() (*api.Client, error) {
	if apiClient != nil {
		return apiClient, nil
	} else {
		var err error

		apiClient, err = api.NewClient(
			conf.TUIConfig.Server.Endpoint,
			conf.TUIConfig.Server.Token,
			internal.CLI_VERSION,
		)

		return apiClient, err
	}
}
