package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/nico-mayer/huectl-cli/config"
	"github.com/nico-mayer/huectl-cli/utils"
	"github.com/reujab/wallpaper"
	"github.com/urfave/cli/v3"
)

type task struct {
	name string
	run  func(themeInfo *utils.ThemeInfo) error
}

func Set() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "sets the specified theme",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:  "Theme",
				Value: "tokyo-night-storm",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			themeName := c.StringArg("Theme")
			themes := utils.FindAvailableThemes()

			valid := slices.Contains(themes, themeName)
			if valid {
				themeInfo, err := utils.SetTheme(themeName)

				if err != nil || themeInfo == nil {
					return err
				}

				tasks := []task{
					{name: "Change zed theme", run: changeZedTheme},
					{name: "Change ghostty theme", run: changeGhosttyTheme},
					{name: "Set wallpaper", run: setWallpaper},
				}

				if runtime.GOOS == "darwin" {
					tasks = append(tasks, task{name: "Set macOS theme", run: setMacOSTheme})
				}

				var wg sync.WaitGroup
				errCh := make(chan error, len(tasks))

				for _, t := range tasks {
					wg.Add(1)
					go func(t task) {
						defer wg.Done()
						err := t.run(themeInfo)
						if err != nil {
							errCh <- fmt.Errorf("%s: %v", t.name, err)
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

			return fmt.Errorf("Invalid theme name")
		},
	}
}

func changeZedTheme(themeInfo *utils.ThemeInfo) error {
	zedSettingsPath := filepath.Join(os.Getenv("HOME"), ".config", "zed", "settings.json")
	data, err := os.ReadFile(zedSettingsPath)
	if err != nil {
		return err
	}

	content := string(data)
	re := regexp.MustCompile(`("theme"\s*:\s*")([^"]*)(")`)
	if !re.MatchString(content) {
		return fmt.Errorf("theme setting not found in %s", zedSettingsPath)
	}

	updatedSettings := re.ReplaceAllString(content, `${1}`+themeInfo.Name+`${3}`)

	err = os.WriteFile(zedSettingsPath, []byte(updatedSettings), 0644)
	if err != nil {
		return err
	}

	return nil
}

func changeGhosttyTheme(themeInfo *utils.ThemeInfo) error {
	ghosttyConfigPath := filepath.Join(os.Getenv("HOME"), ".config", "ghostty", "config.ghostty")

	if len(themeInfo.GhosstyThemeName) == 0 {
		themeInfo.GhosstyThemeName = themeInfo.Name
	}

	data, err := os.ReadFile(ghosttyConfigPath)
	if err != nil {
		return err
	}
	content := string(data)

	re := regexp.MustCompile(`(theme\s*=\s*")([^"]*)(")`)
	if !re.MatchString(content) {
		return fmt.Errorf("theme setting not found in %s", ghosttyConfigPath)
	}

	updatedSettings := re.ReplaceAllString(content, `${1}`+themeInfo.GhosstyThemeName+`${3}`)

	err = os.WriteFile(ghosttyConfigPath, []byte(updatedSettings), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("pkill", "-SIGUSR2", "ghostty")
	_ = cmd.Run()

	return nil
}

func setWallpaper(themeInfo *utils.ThemeInfo) error {
	if len(themeInfo.WallpaperSources) == 0 {
		log.Printf("no wallpaper sources for theme %s", themeInfo.Name)
		return nil
	}

	var supportedFileTypes []string = []string{"png", "jpeg", "jpg", "heic"}

	var validWallpaperPaths []string

	for _, source := range themeInfo.WallpaperSources {
		folderPath := filepath.Join(config.WallpaperDir(), source)

		entires, err := os.ReadDir(folderPath)
		if err != nil {
			continue
		}

		for _, entrie := range entires {
			if !entrie.IsDir() {
				wallpaperPath := filepath.Join(folderPath, entrie.Name())
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
	}

	if len(validWallpaperPaths) > 0 {
		wallpaper.SetFromFile(utils.RandomElement(validWallpaperPaths))
	}

	return nil
}

func setMacOSTheme(themeInfo *utils.ThemeInfo) error {
	mode := strings.ToLower(themeInfo.Appearance)

	var script string

	switch mode {
	case "dark":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to true`
	case "light":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to false`
	default:
		return fmt.Errorf("invalid mode: %s (use 'dark' or 'light')", mode)
	}

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set theme: %v, output: %s", err, string(output))
	}

	return nil
}
