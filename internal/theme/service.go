package theme

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/integrations"
	"github.com/nico-mayer/themectl-cli/internal/model"
	"github.com/nico-mayer/themectl-cli/internal/random"
)

func ListAll() ([]string, error) {
	cfg, _ := config.Get()
	themeDir := cfg.Paths.ThemesDir
	log.Debug("listing available themes", "theme_dir", themeDir)

	var themes []string

	dir, err := os.ReadDir(themeDir)
	if err != nil {
		return nil, fmt.Errorf("read theme directory %q: %w", themeDir, err)
	}

	for _, entry := range dir {
		if entry.IsDir() {
			if !strings.HasPrefix(entry.Name(), "_") {
				themes = append(themes, entry.Name())
			}
		}
	}

	log.Debug("listed available themes", "count", len(themes))
	return themes, nil
}

func Set(themeName string) (model.ThemeInfo, error) {
	cfg, _ := config.Get()
	themes, err := ListAll()
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("list available themes: %w", err)
	}

	valid := slices.Contains(themes, themeName)
	if !valid {
		return model.ThemeInfo{}, fmt.Errorf("theme %q is not available", themeName)
	}

	log.Debug("copying theme into current directory", "theme", themeName)
	themeInfo, err := copyToCurrent(themeName)
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("copy theme %q to current directory: %w", themeName, err)
	}

	tasks := slices.DeleteFunc(integrations.All(), func(i integrations.Integration) bool {
		var shouldBeRemoved bool = !slices.Contains(cfg.Settings.Integrations, i.Name())
		if !shouldBeRemoved {
			log.Debug("loading", "integration", i.Name())
		}
		return shouldBeRemoved
	})

	var wg sync.WaitGroup
	errCh := make(chan error, len(tasks))
	log.Info("start integrations...")
	for _, t := range tasks {
		wg.Add(1)
		go func(t integrations.Integration) {
			defer wg.Done()

			if err := t.Apply(themeInfo); err != nil {
				errCh <- fmt.Errorf("%s: %w", t.Name(), err)
			}
		}(t)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return model.ThemeInfo{}, errors.Join(errs...)
	}

	return themeInfo, nil
}

func SetRandom(appearance string) (model.ThemeInfo, error) {
	allThemeInfos, err := loadAllInfos()
	if err != nil {
		return model.ThemeInfo{}, err
	}

	curentTheme, err := Current()
	if err != nil {
		log.Warn("no current theme found")
	}

	filtered := slices.DeleteFunc(slices.Clone(allThemeInfos), func(info model.ThemeInfo) bool {
		if info.Name == curentTheme.Name {
			return true
		}
		if appearance != "" {
			return !strings.EqualFold(info.Appearance, appearance)
		}
		return false
	})

	if len(filtered) == 0 {
		return model.ThemeInfo{}, fmt.Errorf("no matching themes available")
	}

	selectedTheme := random.Element(filtered)

	return Set(selectedTheme.Name)
}

func Current() (model.ThemeInfo, error) {
	log.Debug("loading current theme info")
	info, err := loadInfo("_current")
	if err != nil {
		return model.ThemeInfo{}, fmt.Errorf("load current theme info: %w", err)
	}

	log.Debug("loaded current theme info", "theme", info.Name)
	return info, nil
}
