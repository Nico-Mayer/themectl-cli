package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/theme"
	"github.com/urfave/cli/v3"
)

func List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list all available themes",
		Action: func(ctx context.Context, c *cli.Command) error {
			log.Debug("listing available themes")

			themes, err := theme.ListAll()
			if err != nil {
				log.Error("failed to list themes", "error", err)
				return fmt.Errorf("list themes: %w", err)
			}

			log.Info("themes loaded", "count", len(themes))

			for _, theme := range themes {
				fmt.Println(theme)
			}

			return nil
		},
	}
}
