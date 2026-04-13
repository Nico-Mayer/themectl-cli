package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/fs"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

func loadInfo(name string) (model.ThemeInfo, error) {
	cfg, _ := config.Get()

	infoFilePath := filepath.Join(cfg.ThemesDir(), name, "info.json")

	log.Debug("loading theme info", "theme", name, "path", infoFilePath)

	data, err := os.ReadFile(filepath.Join(infoFilePath))
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("read theme info file %q: %w", infoFilePath, err)
	}

	var themeInfo model.ThemeInfo

	err = json.Unmarshal(data, &themeInfo)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("parse theme info file %q: %w", infoFilePath, err)
	}

	log.Debug("loaded theme info", "theme", themeInfo.Name, "appearance", themeInfo.Appearance)

	return themeInfo, nil
}

func copyToCurrent(themeName string) (model.ThemeInfo, error) {
	cfg, _ := config.Get()

	srcDir := filepath.Join(cfg.ThemesDir(), themeName)
	targetDir := filepath.Join(cfg.ThemesDir(), "_current")

	log.Debug("copying theme files", "theme", themeName, "source", srcDir, "target", targetDir)

	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("stat theme directory %q: %w", srcDir, err)
	}
	if !srcInfo.IsDir() {
		return model.ThemeInfo{}, fmt.Errorf("theme path %q is not a directory", srcDir)
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return model.ThemeInfo{}, fmt.Errorf("create current theme directory %q: %w", targetDir, err)
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("read current theme directory %q: %w", targetDir, err)
	}
	for _, entry := range entries {
		entryPath := filepath.Join(targetDir, entry.Name())
		log.Debug("removing current theme entry", "path", entryPath)
		if err := os.RemoveAll(entryPath); err != nil {
			return model.ThemeInfo{}, fmt.Errorf("clear current theme entry %q: %w", entryPath, err)
		}
	}

	err = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		dst := filepath.Join(targetDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}

		// if d.Type()&os.ModeSymlink != 0 {
		// 	log.Warn("skipping symlink while copying theme", "path", path)
		// 	return nil
		// }

		log.Debug("copying theme file", "source", path, "target", dst)
		if err := fs.CopyFile(path, dst); err != nil {
			return fmt.Errorf("copy file %q to %q: %w", path, dst, err)
		}

		return nil
	})
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("copy theme directory %q to %q: %w", srcDir, targetDir, err)
	}

	themeInfo, err := Current()
	if err != nil {
		return model.ThemeInfo{}, err
	}

	log.Info("theme files copied", "theme", themeInfo.Name, "appearance", themeInfo.Appearance)
	return themeInfo, nil
}
