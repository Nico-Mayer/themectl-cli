package cmd

import (
	"context"
	"fmt"

	"github.com/nico-mayer/huectl-cli/utils"
	"github.com/urfave/cli/v3"
)

func List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list all available themes",
		Action: func(ctx context.Context, c *cli.Command) error {
			themes := utils.FindAvailableThemes()

			for _, theme := range themes {
				fmt.Println(theme)
			}

			return nil
		},
	}
}
