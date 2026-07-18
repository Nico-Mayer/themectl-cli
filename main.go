package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/Nico-Mayer/themectl/internal/cli"
	"github.com/Nico-Mayer/themectl/internal/config"
	"github.com/Nico-Mayer/themectl/internal/engine"
	"github.com/Nico-Mayer/themectl/internal/integration"
	"github.com/Nico-Mayer/themectl/internal/theme"
	"github.com/charmbracelet/log"
)

func main() {
	root := defaultRoot()
	cfg, err := config.Load(root)
	if err != nil {
		log.Fatal(err)
	}

	store := theme.NewStore(os.DirFS(cfg.ThemesDir()))
	engine := engine.New(integration.Enabled(cfg))

	app := cli.New(cfg, store, engine)
	app.Version = version()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "unknown"
}

func defaultRoot() string {
	userHome, _ := os.UserHomeDir()
	return filepath.Join(userHome, ".config", "themectl")
}
