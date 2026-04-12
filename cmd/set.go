package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/huectl-cli/internal/theme"
	"github.com/urfave/cli/v3"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "sets the specified theme",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:  "Theme",
				Value: "tokyo-night-storm",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			themeName := c.StringArg("Theme")

			log.Info("setting theme", "theme", themeName)

			if err := theme.Set(themeName); err != nil {
				log.Error("failed to set theme", "theme", themeName, "err", err)
				return fmt.Errorf("set theme %q: %w", themeName, err)
			}

			log.Info("theme set successfully", "theme", themeName)
			return nil
		},
	}
}
