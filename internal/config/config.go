package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Paths
	Settings
}

var (
	instance Config
	once     sync.Once
	initErr  error
)

func Get() (Config, error) {
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			initErr = fmt.Errorf("resolve user home dir: %w", err)
			return
		}

		paths := NewPaths(filepath.Join(home, ".config", "themectl"))
		settings, err := LoadSettings(paths.SettingsPath)
		if err != nil {
			initErr = err
			return
		}
		instance = Config{Paths: paths, Settings: settings}
	})

	return instance, initErr
}
