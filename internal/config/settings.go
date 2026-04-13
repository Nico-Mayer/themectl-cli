package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Settings struct {
	Integrations []string `json:"integrations"`
	DefaultTheme string   `json:"default-theme,omitempty"`
}

func DefaultSettings() Settings {
	return Settings{
		Integrations: []string{"ghostty", "zed", "system-theme", "wallpaper", "yazi", "eza"},
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
