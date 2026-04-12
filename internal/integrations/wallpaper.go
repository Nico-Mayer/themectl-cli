package integrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/fs"
	"github.com/nico-mayer/themectl-cli/internal/model"
	"github.com/nico-mayer/themectl-cli/internal/random"
	"github.com/reujab/wallpaper"
)

type Wallpaper struct{}

func (Wallpaper) Name() string {
	return "wallpaper"
}

func (i Wallpaper) Apply(themeInfo model.ThemeInfo) error {
	cfg, _ := config.Get()
	logger := integrationLogger(i)

	sourceDirs := collectWallpaperSourceDirs(cfg, themeInfo, logger)
	logger.Debug("collecting wallpaper candidates", "sources", len(sourceDirs))

	candidates := collectWallpaperCandidates(sourceDirs, logger)
	if len(candidates) == 0 {
		return fmt.Errorf("failed to set wallpaper for theme %q: no supported wallpaper files found", themeInfo.Name)
	}

	current, err := wallpaper.Get()
	if err != nil {
		logger.Warn("no current wallpaper found")
	}

	selected := pickWallpaper(candidates, current, logger)

	logger.Info("setting wallpaper", "selected", selected, "candidates", len(candidates))
	if err := wallpaper.SetFromFile(selected); err != nil {
		return fmt.Errorf("failed to set wallpaper for theme %q from %q: %w", themeInfo.Name, selected, err)
	}

	logger.Info("wallpaper updated successfully", "selected", selected)
	return nil
}

func collectWallpaperSourceDirs(cfg config.Config, themeInfo model.ThemeInfo, logger log.Logger) []string {
	dirs := []string{
		filepath.Join(cfg.ThemesDir(), themeInfo.Name, "wallpaper"),
	}

	for _, source := range themeInfo.WallpaperSources {
		themePath := filepath.Join(cfg.ThemesDir(), source, "wallpaper")
		exists := fs.Exists(themePath)
		if fs.Exists(themePath) {
			dirs = append(dirs, themePath)
		} else {
			logger.Debug("skipping wallpaper source: path does not exist or is not accessible", "path", themePath)
		}

		sourcePath := filepath.Join(cfg.WallpaperSourcesDir(), source)
		exists = fs.Exists(sourcePath)
		if exists {
			dirs = append(dirs, sourcePath)
		} else {
			logger.Debug("skipping wallpaper source: path does not exist or is not accessible", "path", sourcePath)
		}
	}

	return dirs
}

func collectWallpaperCandidates(sourceDirs []string, logger log.Logger) []string {
	supported := map[string]struct{}{
		".png":  {},
		".jpeg": {},
		".jpg":  {},
		".heic": {},
	}

	candidates := make([]string, 0)

	for _, dir := range sourceDirs {
		entries, err := os.ReadDir(dir)

		if err != nil {
			logger.Warn("skipping wallpaper source", "path", dir, "err", err)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			path := filepath.Join(dir, entry.Name())
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if _, ok := supported[ext]; ok {
				candidates = append(candidates, path)
			}
		}
	}

	return candidates
}

func pickWallpaper(candidates []string, current string, logger log.Logger) string {
	if len(candidates) == 1 {
		return candidates[0]
	}

	selected := random.Element(candidates)
	for selected == current {
		logger.Info("reselect wallpaper because it is already set")
		selected = random.Element(candidates)
	}
	return selected
}
