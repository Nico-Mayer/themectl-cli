package integration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nico-Mayer/themectl/internal/config"
	"github.com/Nico-Mayer/themectl/internal/theme"
)

type Eza struct {
	Cfg config.Config
}

func (Eza) Name() string {
	return "eza"
}

func (e Eza) Apply(t theme.Resolved) error {
	sourceFile := filepath.Join(e.Cfg.CurrentDir(), "eza.yml")
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve user home dir: %w", err)
	}
	targetFile := filepath.Join(userHomePath, ".config", "eza", "theme.yml")

	return symlink(sourceFile, targetFile)
}

func (e Eza) Check() error {
	return checkFileExists("eza theme file", filepath.Join(e.Cfg.CurrentDir(), "eza.yml"))
}
