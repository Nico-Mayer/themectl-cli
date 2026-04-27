package integrations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

type Nvim struct{}

func init() {
	Register(Nvim{})
}

func (Nvim) Name() string {
	return "nvim"
}

func (i Nvim) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)

	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve user home dir: %w", err)
	}

	srcPath := filepath.Join(cfg.Paths.CurrentThemeDir, "nvim.lua")
	targetDirPath := filepath.Join(userHomePath, ".config", "nvim", "plugin")
	targetFilePath := filepath.Join(targetDirPath, "99_theme.lua")

	logger.Debug("linking theme", "src", srcPath, "dst", targetFilePath)

	if err := os.MkdirAll(targetDirPath, 0o755); err != nil {
		return fmt.Errorf("create nvim config dir %q: %w", targetDirPath, err)
	}

	if err := os.Symlink(srcPath, targetFilePath); err != nil {
		if errors.Is(err, os.ErrExist) {
			logger.Info("symlink exists, skipping")
			return nil
		}
		return fmt.Errorf("create nvim symlink %q -> %q: %w", targetFilePath, srcPath, err)
	}

	logger.Info("applied", "theme", themeInfo.Name)
	return nil
}
