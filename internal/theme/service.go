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
)

func ListAll() ([]string, error) {
	cfg, _ := config.Get()
	themeDir := cfg.ThemesDir()
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

func Set(themeName string) error {
	themes, err := ListAll()
	if err != nil {
		return fmt.Errorf("list available themes: %w", err)
	}

	valid := slices.Contains(themes, themeName)
	if !valid {
		return fmt.Errorf("theme %q is not available", themeName)
	}

	log.Debug("copying theme into current directory", "theme", themeName)
	themeInfo, err := copyToCurrent(themeName)
	if err != nil {
		return fmt.Errorf("copy theme %q to current directory: %w", themeName, err)
	}

	tasks := []integrations.Integration{
		integrations.Zed{},
		integrations.Ghostty{},
		integrations.SystemTheme{},
		integrations.Wallpaper{},
		integrations.Yazi{},
		integrations.Eza{},
	}

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
		return errors.Join(errs...)
	}

	return nil
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
