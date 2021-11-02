package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var importCommand = &cli.Command{
	Name:  "import",
	Usage: "import GitHub labels",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "repo",
			Usage:    "the name of the GitHub repo to import the labels to",
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
		input := context.Args().Get(0)

		ctx := context.Context

		var inFile *os.File
		var err error

		if input == "-" {
			inFile = os.Stdin
		} else {
			inFile, err = os.OpenFile(input, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
		}
		defer inFile.Close()

		data, err := ioutil.ReadAll(inFile)
		if err != nil {
			return err
		}

		labels := []Label{}
		err = json.Unmarshal(data, &labels)
		if err != nil {
			return err
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		for i := range labels {
			label := labels[i]

			_, resp, err := client.Issues.GetLabel(ctx, org, repo, label.Name)
			if err != nil {
				if resp.StatusCode != http.StatusNotFound {
					return err
				}
				_, _, createErr := client.Issues.CreateLabel(ctx, org, repo, &github.Label{
					Name:        &label.Name,
					Color:       &label.Color,
					Description: label.Description,
					Default:     label.Default,
				})
				if createErr != nil {
					return createErr
				}
				fmt.Printf("created label %s\n", label.Name)
			}
			// TODO: shall we update in the future?
		}

		return nil
	},
}
