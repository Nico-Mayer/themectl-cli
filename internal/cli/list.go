package cli

import (
	"context"
	"fmt"

	"github.com/nico-mayer/themectl-cli/internal/theme"
	"github.com/urfave/cli/v3"
)

func listCmd(store *theme.Store) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List all available themes",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "light",
				Aliases: []string{"l"},
				Usage:   "only list light themes",
			},
			&cli.BoolFlag{
				Name:    "dark",
				Aliases: []string{"d"},
				Usage:   "only list dark themes",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			light := c.Bool("light")
			dark := c.Bool("dark")

			if light && dark {
				return fmt.Errorf("cannot use --light and --dark together")
			}

			if light {
				return listByAppearance(store, theme.Light)
			}
			if dark {
				return listByAppearance(store, theme.Dark)
			}

			all, err := store.ListAll()
			if err != nil {
				return err
			}

			for _, t := range all {
				fmt.Println(t)
			}

			return nil
		},
	}
}

func listByAppearance(store *theme.Store, a theme.Appearance) error {
	all, err := store.ListAllByAppearance(a)
	if err != nil {
		return err
	}

	for _, t := range all {
		fmt.Println(t.ID())
	}
	return nil
}
