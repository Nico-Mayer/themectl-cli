package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/integrations"
	"github.com/nico-mayer/themectl-cli/internal/theme"
	"github.com/reujab/wallpaper"
	"github.com/urfave/cli/v3"
)

func Wallpaper() *cli.Command {
	return &cli.Command{
		Name:    "wallpaper",
		Aliases: []string{"wall"},
		Usage:   "set wallpaper",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "random",
				Aliases: []string{"r"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			randomFlag := c.Bool("random")
			if randomFlag {
				themeInfo, err := theme.Current()
				if err != nil {
					log.Error("failed to get current theme", "error", err)
					return fmt.Errorf("get current theme: %w", err)
				}

				err = integrations.Wallpaper{}.Apply(themeInfo)
				if err != nil {
					themeName := themeInfo.Name
					log.Error("failed to set wallpaper", "theme", themeName, "err", err)
					return fmt.Errorf("set wallpaper %q: %w", themeName, err)
				}
				return nil
			}

			currentWall, err := wallpaper.Get()
			if err != nil {
				return err
			}
			fmt.Println(currentWall)

			return nil
		},
	}
}
