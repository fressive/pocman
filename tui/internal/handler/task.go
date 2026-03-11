package handler

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/tui/internal/api"
	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
)

func NewTask(ctx context.Context, c *cli.Command) error {
	client, err := api.GetClient()

	if err != nil {
		return err
	}

	var vuln model.Vuln

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[model.Vuln]().
				Title("Vulnerability").
				OptionsFunc(func() []huh.Option[model.Vuln] {
					vulns, err := client.ListVulns(ctx)

					if err != nil {
						return nil
					}

					return lo.Map(vulns, func(v model.Vuln, _ int) huh.Option[model.Vuln] {
						return huh.NewOption(fmt.Sprintf("#%d | %s | %s", v.ID, v.Code, v.Title), v)
					})
				}, &vuln).
				Filtering(true),
		),
	).Run()

	if err != nil {
		return err
	}

	return nil
}
