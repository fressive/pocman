package flows

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
)

type GenerateVulnInput struct {
}

type GenerateVulnOutput struct {
}

var GenerateVulnFlow *core.Flow[GenerateVulnInput, *GenerateVulnOutput, struct{}]

func InitVulnFlows(g *genkit.Genkit) {
	GenerateVulnFlow = genkit.DefineFlow(g, "generate_vuln", func(ctx context.Context, input GenerateVulnInput) (*GenerateVulnOutput, error) {
		genkit.GenerateData[GenerateVulnOutput](ctx, g,
			ai.WithPrompt(""),
			ai.WithDocs(&ai.Document{}),
		)

		return nil, nil
	})
}
