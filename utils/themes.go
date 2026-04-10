package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nico-mayer/huectl-cli/config"
)

type ThemeInfo struct {
	Name             string   `json:"name"`
	Appearance       string   `json:"appearance"`
	GhosstyThemeName string   `json:"ghostty-theme-name"`
	WallpaperSources []string `json:"wallpaper-sources"`
}

func FindAvailableThemes() []string {
	var themes []string

	dir, err := os.ReadDir(config.ThemeDir())
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range dir {
		if entry.IsDir() {
			if entry.Name() != "_current" {
				themes = append(themes, entry.Name())
			}
		}
	}

	return themes
}

func SetTheme(name string) (*ThemeInfo, error) {
	srcDir := filepath.Join(config.ThemeDir(), name)
	targetDir := filepath.Join(config.ThemeDir(), "_current")

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read target dir: %w", err)
	}
	for _, entry := range entries {
		entryPath := filepath.Join(targetDir, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return nil, fmt.Errorf("failed to delete %s: %w", entryPath, err)
		}
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create target dir: %w", err)
	}

	srcEntries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source dir: %w", err)
	}
	for _, entry := range srcEntries {
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(targetDir, entry.Name())
		if err := copyFile(src, dst); err != nil {
			return nil, fmt.Errorf("failed to copy %s: %w", entry.Name(), err)
		}
	}

	themeInfo := GetCurrentThemeInfo()

	return &themeInfo, nil
}

func GetCurrentThemeInfo() ThemeInfo {
	return getThemeInfo("_current")
}

func getThemeInfo(name string) ThemeInfo {
	infoFilePath := filepath.Join(config.ThemeDir(), name, "info.json")

	data, err := os.ReadFile(filepath.Join(infoFilePath))
	if err != nil {
		log.Fatal("failed to load theme data")
	}

	var themeInfo ThemeInfo

	err = json.Unmarshal(data, &themeInfo)
	if err != nil {
		log.Fatal("failed to unmarshal theme data")
	}

	return themeInfo
}
