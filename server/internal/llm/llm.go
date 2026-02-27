package llm

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai"
	"github.com/firebase/genkit/go/plugins/compat_oai/anthropic"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/fressive/pocman/server/internal/conf"
	"github.com/openai/openai-go/option"
)

// Define the input structure for the tool
type StringInput struct {
	RawString string `json:"raw_string" jsonschema_description:"The raw string to input."`
}

func NewAgent() (*genkit.Genkit, error) {
	var g *genkit.Genkit

	if llmConfig := conf.ServerConfig.LLM; llmConfig != nil {
		ctx := context.Background()

		switch llmConfig.Provider {
		case "anthropic":
			g = genkit.Init(ctx, genkit.WithPlugins(&anthropic.Anthropic{
				Opts: []option.RequestOption{
					option.WithAPIKey(llmConfig.APIKey),
					option.WithBaseURL(llmConfig.Endpoint),
				},
			}))
		case "googlegenai":
			g = genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{
				APIKey: llmConfig.APIKey,
			}))
		case "custom":
			g = genkit.Init(ctx, genkit.WithPlugins(&compat_oai.OpenAICompatible{
				APIKey:   llmConfig.APIKey,
				BaseURL:  llmConfig.Endpoint,
				Provider: "custom",
			}), genkit.WithDefaultModel(llmConfig.Model))

		default:
			return nil, fmt.Errorf("config error: unknown llm provider type %s", llmConfig.Provider)
		}

		return g, nil

	} else {
		return nil, fmt.Errorf("config error: llm config not found")
	}
}
