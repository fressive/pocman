package tool

import (
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

type DeployDockerComposeInput struct {
	DockerComposeRaw string `json:"docker_compose_raw" jsonschema_description:"The content of generated docker-compose.yml file."`
}

func DeployDockerCompose(ctx *ai.ToolContext, input DeployDockerComposeInput) (string, error) {
	fmt.Println(input.DockerComposeRaw)
	compose, _ := compose.NewDockerComposeWith(compose.WithStackReaders(strings.NewReader(input.DockerComposeRaw)))

	err := compose.Up(ctx)
	if err != nil {
		return "Fail to deploy", err
	}

	return "Deploy successfully", nil
}
