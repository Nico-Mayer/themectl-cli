package cli

import (
	"context"

	"github.com/nico-mayer/themectl-cli/internal/theme"
	"github.com/urfave/cli/v3"
)

func listCommand(store *theme.Store) *cli.Command {
	return &cli.Command{
		Name: "list",
		Action: func(ctx context.Context, c *cli.Command) error {

			return nil
		},
	}
}
