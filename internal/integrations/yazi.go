package integrations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

type Yazi struct{}

func (Yazi) Name() string {
	return "yazi"
}

func (i Yazi) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("unable to resolve user home dir: %w", err)
	}

	cfg, _ := config.Get()
	yaziFlavorFilePath := filepath.Join(cfg.CurrentThemeDir(), "yazi-flavor.toml")
	targetDir := filepath.Join(userHomeDir, ".config", "yazi", "flavors", "themectl.yazi")
	linkPath := filepath.Join(targetDir, "flavor.toml")

	err = os.MkdirAll(targetDir, 493)
	if err != nil {
		return err
	}

	if err := os.Symlink(yaziFlavorFilePath, linkPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			logger.Debug("yazi flavor symlink already exists, skipping", "path", linkPath)
			return nil
		}
		return fmt.Errorf("create yazi symlink %q -> %q: %w", linkPath, yaziFlavorFilePath, err)
	}

	return nil
}
