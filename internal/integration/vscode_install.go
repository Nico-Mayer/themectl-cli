package integration

import (
	"log/slog"

	"github.com/Nico-Mayer/themectl/internal/cache"
)

type codeInstaller struct {
	cache   cache.Cache
	install func(string) error
}

func (c codeInstaller) Ensure(id string) error {
	if c.cache.Fresh(id, installCheckTTL) {
		slog.Debug("vscode extension recently checked, skipping", "id", id)
		return nil
	}

	slog.Debug("installing vscode extension", "id", id)
	if err := c.install(id); err != nil {
		return err
	}
	slog.Info("vscode extension ensured", "id", id)

	return c.cache.Put(id, []byte("installed"))
}
