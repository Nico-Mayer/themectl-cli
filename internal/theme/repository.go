package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

var themeCache map[string]model.ThemeInfo = make(map[string]model.ThemeInfo)

func loadInfo(name string) (model.ThemeInfo, error) {
	themeInfo, ok := themeCache[name]
	if ok {
		log.Debug("loaded from cache", "theme", themeInfo.Name)
		return themeInfo, nil
	}

	cfg, _ := config.Get()

	infoFilePath := filepath.Join(cfg.Paths.ThemesDir, name, "info.json")

	log.Debug("loading theme info", "theme", name, "path", infoFilePath)

	data, err := os.ReadFile(filepath.Join(infoFilePath))
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("read theme info file %q: %w", infoFilePath, err)
	}

	err = json.Unmarshal(data, &themeInfo)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("parse theme info file %q: %w", infoFilePath, err)
	}

	log.Debug("loaded theme info", "theme", themeInfo.Name, "appearance", themeInfo.Appearance)
	themeCache[name] = themeInfo

	return themeInfo, nil
}

func loadAllInfos() ([]model.ThemeInfo, error) {
	cfg, _ := config.Get()

	dir, err := os.ReadDir(cfg.ThemesDir)
	if err != nil {
		return []model.ThemeInfo{}, err
	}

	var filteredEntries []string
	for _, d := range dir {
		if d.IsDir() && !strings.HasPrefix(d.Name(), "_") {
			filteredEntries = append(filteredEntries, d.Name())
		}
	}

	var wg sync.WaitGroup
	resultChan := make(chan model.ThemeInfo, len(filteredEntries))

	for _, e := range filteredEntries {
		wg.Add(1)
		go func(themeName string) {
			defer wg.Done()
			info, err := loadInfo(themeName)
			if err != nil {
				log.Warn("skipping theme", "theme", themeName, "err", err)
				return
			}

			resultChan <- info
		}(e)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var themeInfos []model.ThemeInfo
	for info := range resultChan {
		themeInfos = append(themeInfos, info)
	}

	if len(themeInfos) == 0 {
		return []model.ThemeInfo{}, fmt.Errorf("no theme could be loaded")
	}
	return themeInfos, nil
}

func symLinkToCurrent(themeName string) (model.ThemeInfo, error) {
	cfg, _ := config.Get()

	srcDir := filepath.Join(cfg.Paths.ThemesDir, themeName)
	targetDir := filepath.Join(cfg.Paths.ThemesDir, "_current")

	log.Debug("linking current theme", "theme", themeName, "source", srcDir, "target", targetDir)

	themeInfo, err := loadInfo(themeName)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("load theme info for %q: %w", themeName, err)
	}

	if err := os.RemoveAll(targetDir); err != nil {
		return model.ThemeInfo{}, fmt.Errorf("remove existing current theme link %q: %w", targetDir, err)
	}

	if err := os.Symlink(srcDir, targetDir); err != nil {
		return model.ThemeInfo{}, fmt.Errorf("create symlink %q -> %q: %w", targetDir, srcDir, err)
	}

	log.Debug("linked current theme", "theme", themeName)
	return themeInfo, nil
}
