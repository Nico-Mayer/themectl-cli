package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/huectl-cli/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "huectl",
		Usage: "my theme switcher cli stuffi",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if c.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			cmd.List(),
			cmd.Set(),
			cmd.Current(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(fmt.Errorf("huectl failed: %w", err))
	}
}
