package integrations

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/model"
	"github.com/nico-mayer/themectl-cli/internal/random"
	"github.com/reujab/wallpaper"
)

func SetWallpaper(themeInfo model.ThemeInfo) error {
	if len(themeInfo.WallpaperSources) == 0 {
		log.Info("no wallpaper sources configured")
		return nil
	}
	currWallpaper, err := wallpaper.Get()
	if err != nil {
		log.Warn("no current wallpaper found")
	}

	supportedFileTypes := []string{"png", "jpeg", "jpg", "heic"}
	validWallpaperPaths := make([]string, 0)

	log.Debug("collecting wallpaper candidates", "sources", len(themeInfo.WallpaperSources))

	for _, source := range themeInfo.WallpaperSources {
		folderPath := filepath.Join(config.WallpaperDir(), source)

		entries, err := os.ReadDir(folderPath)
		if err != nil {
			log.Warn("skipping wallpaper source", "source", source, "path", folderPath, "err", err)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			wallpaperPath := filepath.Join(folderPath, entry.Name())
			pathSubstring := strings.Split(wallpaperPath, ".")

			var fileType string

			if len(pathSubstring) > 0 {
				fileType = strings.ToLower(pathSubstring[len(pathSubstring)-1])
			}

			if slices.Contains(supportedFileTypes, fileType) {
				validWallpaperPaths = append(validWallpaperPaths, wallpaperPath)
			}
		}
	}

	if len(validWallpaperPaths) == 0 {
		return fmt.Errorf("failed to set wallpaper for theme %q: no supported wallpaper files found", themeInfo.Name)
	}

	selectedWallpaper := random.Element(validWallpaperPaths)

	for {
		if selectedWallpaper != currWallpaper || len(validWallpaperPaths) <= 1 {
			break
		}
		log.Info("reselect wallpaper because it is already set")
		selectedWallpaper = random.Element(validWallpaperPaths)
	}

	log.Info("setting wallpaper", "selected", selectedWallpaper, "candidates", len(validWallpaperPaths))

	if err := wallpaper.SetFromFile(selectedWallpaper); err != nil {
		return fmt.Errorf("failed to set wallpaper for theme %q from %q: %w", themeInfo.Name, selectedWallpaper, err)
	}

	log.Info("wallpaper updated successfully", "selected", selectedWallpaper)

	return nil
}
