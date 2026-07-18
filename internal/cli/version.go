package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func versionCmd() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "print the themectl version",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println(c.Root().Version)
			return nil
		},
	}
}
