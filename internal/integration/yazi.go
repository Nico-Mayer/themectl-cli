package integration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nico-Mayer/themectl/internal/config"
	"github.com/Nico-Mayer/themectl/internal/theme"
)

type Yazi struct {
	Cfg config.Config
}

func (Yazi) Name() string {
	return "yazi"
}

func (y Yazi) Apply(t theme.Resolved) error {
	sourceFile := filepath.Join(y.Cfg.CurrentDir(), "yazi-flavor.toml")
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve user home dir: %w", err)
	}

	targetFile := filepath.Join(userHomePath, ".config", "yazi", "flavors", "themectl.yazi", "flavor.toml")

	return symlink(sourceFile, targetFile)
}

func (y Yazi) Check() error {
	return checkFileExists("yazi flavor file", filepath.Join(y.Cfg.CurrentDir(), "yazi-flavor.toml"))
}
