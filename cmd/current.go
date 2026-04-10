package cmd

import (
	"context"
	"fmt"

	"github.com/nico-mayer/huectl-cli/utils"
	"github.com/urfave/cli/v3"
)

type ThemeInfo struct {
	Name             string   `json:"name"`
	Appearance       string   `json:"appearance"`
	GhosstyThemeName string   `json:"ghostty-theme-name"`
	WallpaperSources []string `json:"wallpaper-sources"`
}

func Current() *cli.Command {
	return &cli.Command{
		Name:  "current",
		Usage: "get the current active theme",
		Action: func(ctx context.Context, c *cli.Command) error {
			themeInfo := utils.GetCurrentThemeInfo()

			fmt.Println(themeInfo.Name)

			return nil
		},
	}
}
