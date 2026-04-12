package theme

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/huectl-cli/internal/config"
	"github.com/nico-mayer/huectl-cli/internal/integrations"
	"github.com/nico-mayer/huectl-cli/internal/model"
)

type task struct {
	name string
	run  func(themeInfo model.ThemeInfo) error
}

func ListAll() ([]string, error) {
	themeDir := config.ThemeDir()
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

	tasks := []task{
		{name: "change zed theme", run: integrations.ChangeZedTheme},
		{name: "change ghostty theme", run: integrations.ChangeGhosttyTheme},
		{name: "set wallpaper", run: integrations.SetWallpaper},
		{name: "set system theme", run: integrations.SetSystemTheme},
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(tasks))

	for _, t := range tasks {
		wg.Add(1)
		go func(t task) {
			defer wg.Done()

			err := t.run(themeInfo)
			if err != nil {
				errCh <- fmt.Errorf("%s: %w", t.name, err)
				return
			}

			errCh <- nil
		}(t)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
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
