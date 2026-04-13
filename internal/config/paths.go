package config

import "path/filepath"

type Paths struct {
	ConfigDir           string
	ThemesDir           string
	CurrentThemeDir     string
	WallpaperSourcesDir string
	ReceiptsDir         string
	SettingsPath        string
}

func NewPaths(configDir string) Paths {
	return Paths{
		ConfigDir:           configDir,
		ThemesDir:           filepath.Join(configDir, "themes"),
		CurrentThemeDir:     filepath.Join(configDir, "themes", "_current"),
		WallpaperSourcesDir: filepath.Join(configDir, "wallpaper-sources"),
		ReceiptsDir:         filepath.Join(configDir, "receips"),
		SettingsPath:        filepath.Join(configDir, "themectl.json"),
	}
}
