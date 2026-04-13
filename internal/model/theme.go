package model

type ThemeInfo struct {
	Name             string            `json:"name"`
	Appearance       string            `json:"appearance"`
	WallpaperSources []string          `json:"wallpaper-sources"`
	Overrides        map[string]string `json:"overrides"`
}
