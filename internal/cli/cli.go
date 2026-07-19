package cli

import (
	"context"
	"log/slog"
	"os"

	"github.com/Nico-Mayer/themectl/internal/config"
	"github.com/Nico-Mayer/themectl/internal/engine"
	"github.com/Nico-Mayer/themectl/internal/theme"
	"github.com/charmbracelet/log"
	urfaveCli "github.com/urfave/cli/v3"
)

type app struct {
	cfg    config.Config
	store  *theme.Store
	engine *engine.Engine
}

func New(cfg config.Config, store *theme.Store, engine *engine.Engine) *urfaveCli.Command {

	app := app{
		cfg:    cfg,
		store:  store,
		engine: engine,
	}

	return &urfaveCli.Command{
		Name:                  "themectl",
		Usage:                 "Manage and apply themes across your tools",
		EnableShellCompletion: true,
		Flags: []urfaveCli.Flag{
			&urfaveCli.BoolFlag{
				Name:    "verbose",
				Usage:   "Prints more logs to stderr",
				Aliases: []string{"v"},
			},
		},
		Commands: []*urfaveCli.Command{
			app.listCmd(),
			app.setCmd(),
			app.currentCmd(),
			app.wallpaperCmd(),
			app.refreshCmd(),
			app.doctorCmd(),
		},
		Before: func(ctx context.Context, c *urfaveCli.Command) (context.Context, error) {
			level := slog.LevelInfo
			withTime := false
			if c.Bool("verbose") {
				level = slog.LevelDebug
				withTime = true
			}
			handler := log.NewWithOptions(os.Stderr, log.Options{
				Level:           log.Level(level),
				ReportTimestamp: withTime,
				TimeFormat:      "15:04:05",
			})
			slog.SetDefault(slog.New(handler))

			return nil, nil
		},
	}
}
