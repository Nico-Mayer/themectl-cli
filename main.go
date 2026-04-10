package main

import (
	"context"
	"log"
	"os"

	"github.com/nico-mayer/huectl-cli/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "theme-switcher",
		Usage: "my theme switcher cli stuffi",
		Commands: []*cli.Command{
			cmd.List(),
			cmd.Set(),
			cmd.Current(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
