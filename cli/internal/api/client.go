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
	"github.com/fressive/pocman/common/pkg/model"
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

func (c *Client) Do(ctx context.Context, method, path string, body any, res any) error {
	var bodyReader io.Reader
	contentType := "application/json"
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}

	return c.doRequest(ctx, method, path, contentType, bodyReader, res)
}

func (c *Client) DoMultipart(ctx context.Context, method, path, contentType string, body io.Reader, res any) error {
	return c.doRequest(ctx, method, path, contentType, body, res)
}

func (c *Client) doRequest(ctx context.Context, method, path, contentType string, body io.Reader, res any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, body)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "pocman-cli/"+internal.CLI_VERSION)
	req.Header.Set("Authorization", "Bearer "+
		base64.RawURLEncoding.EncodeToString([]byte(c.Token)))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var resModel model.Response[json.RawMessage]
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err := json.NewDecoder(resp.Body).Decode(&resModel); err == nil {
			return APIError{
				API:    path,
				Method: method,
				Code:   resModel.Code,
				Msg:    resModel.Msg,
			}
		}
		return fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&resModel); err != nil {
		return err
	}

	if resModel.Code != 0 {
		return APIError{
			API:    path,
			Method: method,
			Code:   resModel.Code,
			Msg:    resModel.Msg,
		}
	}

	if res != nil && len(resModel.Data) > 0 {
		if err := json.Unmarshal(resModel.Data, res); err != nil {
			return err
		}
	}

	return nil
}
