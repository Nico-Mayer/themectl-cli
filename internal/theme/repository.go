package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/huectl-cli/internal/config"
	"github.com/nico-mayer/huectl-cli/internal/fs"
	"github.com/nico-mayer/huectl-cli/internal/model"
)

func loadInfo(name string) (model.ThemeInfo, error) {
	infoFilePath := filepath.Join(config.ThemeDir(), name, "info.json")

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
	srcDir := filepath.Join(config.ThemeDir(), themeName)
	targetDir := filepath.Join(config.ThemeDir(), "_current")

	log.Info("copying theme files", "theme", themeName, "source", srcDir, "target", targetDir)

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

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return model.ThemeInfo{}, fmt.Errorf("create current theme directory %q: %w", targetDir, err)
	}

	srcEntries, err := os.ReadDir(srcDir)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("read theme directory %q: %w", srcDir, err)
	}
	for _, entry := range srcEntries {
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(targetDir, entry.Name())
		log.Debug("copying theme file", "source", src, "target", dst)
		if err := fs.CopyFile(src, dst); err != nil {
			return model.ThemeInfo{}, fmt.Errorf("copy theme file %q to %q: %w", src, dst, err)
		}
	}

	themeInfo, err := Current()
	if err != nil {
		return model.ThemeInfo{}, err
	}

	log.Info("theme files copied", "theme", themeInfo.Name, "appearance", themeInfo.Appearance)

	return themeInfo, nil
}
