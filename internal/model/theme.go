package model

type ThemeInfo struct {
	Name             string   `json:"name"`
	Appearance       string   `json:"appearance"`
	GhosttyThemeName string   `json:"ghostty-theme-name"`
	WallpaperSources []string `json:"wallpaper-sources"`
}
