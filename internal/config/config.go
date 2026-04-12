package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	configDir           string
	configPath          string
	themesDir           string
	wallpaperSourcesDir string
}

func (c Config) ConfigDir() string           { return c.configDir }
func (c Config) ConfigPath() string          { return c.configPath }
func (c Config) ThemesDir() string           { return c.themesDir }
func (c Config) CurrentThemeDir() string     { return filepath.Join(c.themesDir, "_current") }
func (c Config) WallpaperSourcesDir() string { return c.wallpaperSourcesDir }

var (
	instance Config
	once     sync.Once
	initErr  error
)

func Get() (Config, error) {
	once.Do(func() {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			initErr = fmt.Errorf("unable to resolve user home dir: %w", err)
			return
		}
		userConfigDir := filepath.Join(userHomeDir, ".config")

		configDir := filepath.Join(userConfigDir, "themectl")
		instance = Config{
			configDir:           configDir,
			configPath:          filepath.Join(configDir, "themectl.json"),
			themesDir:           filepath.Join(configDir, "themes"),
			wallpaperSourcesDir: filepath.Join(configDir, "wallpaper-sources"),
		}
	})

	if initErr != nil {
		return Config{}, initErr
	}
	return instance, nil
}
