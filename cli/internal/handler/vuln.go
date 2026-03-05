package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/fressive/pocman/cli/internal/api"
	"github.com/fressive/pocman/cli/internal/util"
	"github.com/fressive/pocman/common/pkg/model"
	"github.com/urfave/cli/v3"
)

func createVuln(ctx context.Context, vuln *model.Vuln) error {
	var documentsText string
	var resourcesText string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Value(&vuln.Title),
			huh.NewInput().
				Title("Code").
				Description("CVE ID, CNVD ID, etc.").
				Value(&vuln.Code),
			huh.NewText().
				Title("Description").
				Value(&vuln.Description),
		),
		huh.NewGroup(
			huh.NewText().
				Title("Documents").
				Description("Text documents will be provided to LLM to reproduce the vulnerability - docker compose file, writeups, etc.\nInput the local path of the resources below and split them with a line break.\nDirectories are supported.").
				Value(&documentsText),
		),
		huh.NewGroup(
			huh.NewText().
				Title("Resources").
				Description("Resources will be provided to agent into the reproduction container directly - app source, exploit scripts, writeups, etc.\nInput the local path of the resources below and split them with a line break.\nDirectories are supported.").
				Value(&resourcesText),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	vuln.Title = strings.TrimSpace(vuln.Title)
	vuln.Code = strings.TrimSpace(vuln.Code)
	vuln.Description = strings.TrimSpace(vuln.Description)
	if vuln.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if vuln.Code == "" {
		return fmt.Errorf("code cannot be empty")
	}

	documents := util.NormalizePathArgs(strings.Split(documentsText, "\n"))
	resources := util.NormalizePathArgs(strings.Split(resourcesText, "\n"))

	client, err := api.GetClient()
	if err != nil {
		return err
	}

	var createdVuln model.Vuln
	err = spinner.New().
		Context(ctx).
		Title("Creating vulnerability...").
		ActionWithErr(func(ctx context.Context) error {
			createdVuln, err = client.CreateVuln(ctx, *vuln)
			return err
		}).
		Run()

	if err != nil {
		return err
	}

	fmt.Printf("pocman: vulnerability created, ID=%d\n", createdVuln.ID)

	vulnID := uint64(createdVuln.ID)

	docFiles, err := util.ExpandUploadPaths(documents)
	if err != nil {
		return fmt.Errorf("resolve document paths: %w", err)
	}
	resFiles, err := util.ExpandUploadPaths(resources)
	if err != nil {
		return fmt.Errorf("resolve resource paths: %w", err)
	}

	if len(docFiles) == 0 && len(resFiles) == 0 {
		return nil
	}

	uploaded := 0
	for _, file := range docFiles {
		err := spinner.New().
			Context(ctx).
			Title(fmt.Sprintf("Uploading document %s", file)).
			ActionWithErr(func(ctx context.Context) error {
				_, err := client.UploadFile(ctx, file, vulnID, model.Document)

				if err == nil {
					uploaded++
				}

				fmt.Printf("pocman: document %s uploaded\n", file)

				return err
			}).
			Run()

		if err != nil {
			return err
		}
	}

	for _, file := range resFiles {
		err := spinner.New().
			Context(ctx).
			Title(fmt.Sprintf("Uploading resource %s", file)).
			ActionWithErr(func(ctx context.Context) error {
				_, err := client.UploadFile(ctx, file, vulnID, model.Resource)

				if err == nil {
					uploaded++
				}

				fmt.Printf("pocman: resource %s uploaded\n", file)

				return err
			}).
			Run()

		if err != nil {
			return err
		}
	}

	return nil
}

func CreateVuln(ctx context.Context, cmd *cli.Command) error {
	return createVuln(ctx, &model.Vuln{Product: &model.Product{}})
}

func CreateVulnFromCVE(ctx context.Context, cmd *cli.Command) error {
	cveCode := cmd.Args().First()

	if cveCode == "" {
		err := huh.NewInput().
			Title("Input the CVE code (e.g. CVE-2026-1000): ").
			Value(&cveCode).
			Run()

		if err != nil {
			return err
		}
	}

	cveCode = strings.TrimSpace(cveCode)
	if !strings.HasPrefix(cveCode, "CVE") || len(strings.Split(cveCode, "-")) != 3 {
		return fmt.Errorf("failed to parse the CVE code, check the format (e.g. CVE-2026-1000)")
	}

	cveClient := util.NewClient()
	cve, err := cveClient.Fetch(ctx, cveCode)
	if err != nil {
		return err
	}

	return createVuln(ctx, &model.Vuln{
		Title:       cve.Containers.CNA.Title,
		Code:        cveCode,
		Description: cve.Containers.CNA.Descriptions[0].Value,
		Product:     &model.Product{},
	})
}
