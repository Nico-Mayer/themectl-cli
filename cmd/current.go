package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/theme"
	"github.com/urfave/cli/v3"
)

func Current() *cli.Command {
	return &cli.Command{
		Name:  "current",
		Usage: "get the current active theme",
		Action: func(ctx context.Context, c *cli.Command) error {
			log.Debug("getting current theme")

			themeInfo, err := theme.Current()
			if err != nil {
				log.Error("failed to get current theme", "error", err)
				return fmt.Errorf("get current theme: %w", err)
			}

			log.Info("current theme loaded", "theme", themeInfo.Name)
			fmt.Println(themeInfo.Name)

			return nil
		},
	}
}
