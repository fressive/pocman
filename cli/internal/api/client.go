package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fressive/pocman/cli/internal"
	"github.com/fressive/pocman/cli/internal/conf"
)

type Client struct {
	Endpoint   string
	Token      string
	HTTPClient *http.Client
}

func NewClient(endpoint, token string) (*Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("API endpoint cannot be empty")
	}

	if token == "" {
		return nil, fmt.Errorf("Token cannot be empty")
	}

	// remove the ending /
	endpoint = strings.TrimSuffix(endpoint, "/")

	return &Client{
		Endpoint: endpoint,
		Token:    token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

var apiClient *Client

func GetClient() (*Client, error) {
	if apiClient != nil {
		return apiClient, nil
	} else {
		var err error

		apiClient, err = NewClient(
			conf.CLIConfig.Server.Endpoint,
			conf.CLIConfig.Server.Token,
		)

		return apiClient, err
	}
}

func (c *Client) Do(ctx context.Context, method, path string, body interface{}, res interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "pocman-cli/"+internal.CLI_VERSION)
	req.Header.Set("Authorization", "Bearer "+
		base64.RawURLEncoding.EncodeToString([]byte(c.Token)))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	if res != nil {
		return json.NewDecoder(resp.Body).Decode(res)
	}

	return nil
}
