package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ghlabels",
		Usage: "A utility for import/exporting GitHub labels",
		Commands: []*cli.Command{
			exportCommand,
			importCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
