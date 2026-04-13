package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/theme"
	"github.com/urfave/cli/v3"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "sets the specified theme",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "theme",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "random",
				Aliases: []string{"r"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			themeName := c.StringArg("theme")
			randomTheme := c.Bool("random")

			if randomTheme && len(themeName) > 0 {
				return fmt.Errorf("cannot use --random with an explicit theme name")
			}

			if randomTheme {
				log.Info("setting random theme")
				return theme.SetRandom()
			}

			if len(themeName) == 0 {
				return fmt.Errorf("no theme provided")
			}

			log.Info("setting", "theme", themeName)
			return theme.Set(themeName)
		},
	}
}
