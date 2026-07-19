package integration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nico-Mayer/themectl/internal/theme"
)

type SymlinkIntegration struct {
	IntegrationName string
	SourceFile      string
	Target          string
}

func (s SymlinkIntegration) Name() string { return s.IntegrationName }

func (s SymlinkIntegration) Apply(theme.Resolved) error {
	return symlink(s.SourceFile, s.Target)
}

func (s SymlinkIntegration) Check() error {
	return checkFileExists(s.IntegrationName+" theme asset", s.SourceFile)
}

func symlink(source, target string) error {
	if _, err := os.Stat(source); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source file does not exist %q", source)
		}
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creating parent folders %q: %w", target, err)
	}

	err := os.Symlink(source, target)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("creating symlink %q -> %q: %w", source, target, err)
	}

	dest, err := os.Readlink(target)
	if err != nil {
		return fmt.Errorf("target is not a symlink %v", err)
	}

	if dest == source {
		return nil
	}

	if err := os.Remove(target); err != nil {
		return fmt.Errorf("removing stale symlink %q: %w", target, err)
	}

	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("recreating symlink %q -> %q: %w", source, target, err)
	}

	return nil
}
