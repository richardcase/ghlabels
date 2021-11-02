package main

import (
	"encoding/json"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var exportCommand = &cli.Command{
	Name:  "export",
	Usage: "export GitHub labels",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "repo",
			Usage:    "the name of the GitHub repo to export the labels from",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "org",
			Usage:    "the name of the GitHub organization for the repo",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "token",
			Usage:    "the GitHub token to use",
			Required: true,
			EnvVars:  []string{"GITHUB_TOKEN"},
		},
	},
	Action: func(context *cli.Context) error {
		repo := context.String("repo")
		org := context.String("org")
		token := context.String("token")
		output := context.Args().Get(0)

		ctx := context.Context

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		labels, _, err := client.Issues.ListLabels(ctx, org, repo, &github.ListOptions{})
		if err != nil {
			return nil
		}

		exportLabels := []Label{}
		for i := range labels {
			label := labels[i]

			exportLabels = append(exportLabels, Label{
				Name:        *label.Name,
				Color:       *label.Color,
				Description: label.Description,
				Default:     label.Default,
			})
		}

		var out *os.File

		if output == "-" {
			out = os.Stdout
		} else {
			out, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			defer out.Close()
		}

		data, err := json.Marshal(exportLabels)
		if err != nil {
			return err
		}

		_, err = out.Write(data)
		if err != nil {
			return err
		}

		return nil
	},
}
