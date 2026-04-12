package config

import (
	"os"
	"path/filepath"
)

func ThemeDir() string {
	return filepath.Join(os.Getenv("HOME"), ".dotfiles", "themes")
}

func WallpaperDir() string {
	return filepath.Join(ThemeDir(), "_wallpapers")
}
