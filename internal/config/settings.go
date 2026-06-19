package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Settings struct {
	Integrations []string          `json:"integrations"`
	DefaultTheme string            `json:"default-theme,omitempty"`
	ConfigPaths  map[string]string `json:"configpaths,omitempty"`
}

func DefaultSettings() Settings {
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = os.Getenv("HOME")
	}

	return Settings{
		Integrations: []string{
			"ghostty",
			"zed",
			"system-theme",
			"wallpaper",
			"yazi",
			"eza",
			"nvim",
			"helix",
		},
		ConfigPaths: map[string]string{
			"ghostty": filepath.Join(userHome, ".config", "ghostty", "config.ghostty"),
			"zed":     zedConfigPath(userHome),
			"helix":   filepath.Join(userHome, ".config", "helix", "config.toml"),
		},
	}
}

func LoadSettings(path string) (Settings, error) {
	defaults := DefaultSettings()

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaults, nil
	}
	if err != nil {
		return Settings{}, fmt.Errorf("read settings %w", err)
	}

	// This merges because unmarshal is overwriting in existing struct
	if err := json.Unmarshal(data, &defaults); err != nil {
		return Settings{}, fmt.Errorf("parse settings: %w", err)
	}

	return defaults, nil
}

func (s Settings) ConfigPathFor(integration string) string {
	if s.ConfigPaths == nil {
		return ""
	}

	path, ok := s.ConfigPaths[integration]
	if !ok {
		return ""
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}

	path = os.ExpandEnv(path)

	if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
		return path
	}

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}

	return path
}

func zedConfigPath(userHome string) string {
	platform := runtime.GOOS

	configHome, _ := os.UserConfigDir()
	if platform == "windows" {
		return filepath.Join(configHome, "zed", "settings.json")
	}

	return filepath.Join(userHome, ".config", "zed", "settings.json")
}
