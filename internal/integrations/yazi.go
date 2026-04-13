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

func init() {
	Register(Yazi{})
}

func (Yazi) Name() string {
	return "yazi"
}

func (i Yazi) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve user home dir: %w", err)
	}

	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	sourcePath := filepath.Join(cfg.Paths.CurrentThemeDir, "yazi-flavor.toml")
	targetDir := filepath.Join(userHomeDir, ".config", "yazi", "flavors", "themectl.yazi")
	linkPath := filepath.Join(targetDir, "flavor.toml")

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create yazi flavor directory %q: %w", targetDir, err)
	}

	if err := os.Symlink(sourcePath, linkPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			logger.Info("symlink exists, skipping")
			return nil
		}
		return fmt.Errorf("create yazi symlink %q -> %q: %w", linkPath, sourcePath, err)
	}

	logger.Info("applied", "theme", themeInfo.Name)
	return nil
}
