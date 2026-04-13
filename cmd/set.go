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
		Commands: []*cli.Command{
			{
				Name:  "random",
				Usage: "sets a random theme",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "light",
						Aliases: []string{"l"},
						Usage:   "only use light themes",
					},
					&cli.BoolFlag{
						Name:    "dark",
						Aliases: []string{"d"},
						Usage:   "only use dark themes",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					light := c.Bool("light")
					dark := c.Bool("dark")

					if light && dark {
						return fmt.Errorf("cannot use --light and --dark together")
					}

					var appearance string
					if light {
						appearance = "light"
					} else if dark {
						appearance = "dark"
					}

					log.Info("setting random theme", "appearance", appearance)
					themeInfo, err := theme.SetRandom(appearance)
					if err != nil {
						return err
					}
					fmt.Println(themeInfo.Name)
					return nil
				},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			themeName := c.StringArg("theme")

			if len(themeName) == 0 {
				return fmt.Errorf("no theme provided")
			}

			log.Info("setting", "theme", themeName)
			themeInfo, err := theme.Set(themeName)
			if err != nil {
				return err
			}
			fmt.Println(themeInfo.Name)
			return nil
		},
	}
}
