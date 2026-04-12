package integrations

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/huectl-cli/internal/model"
)

func ChangeZedTheme(themeInfo model.ThemeInfo) error {
	zedSettingsPath := filepath.Join(os.Getenv("HOME"), ".config", "zed", "settings.json")

	log.Debug("updating Zed theme", "theme", themeInfo.Name, "path", zedSettingsPath)

	data, err := os.ReadFile(zedSettingsPath)
	if err != nil {
		return fmt.Errorf("read Zed settings from %s: %w", zedSettingsPath, err)
	}

	content := string(data)
	re := regexp.MustCompile(`("theme"\s*:\s*")([^"]*)(")`)
	if !re.MatchString(content) {
		return fmt.Errorf("update Zed theme in %s: could not find \"theme\" setting", zedSettingsPath)
	}

	updatedSettings := re.ReplaceAllString(content, `${1}`+themeInfo.Name+`${3}`)

	err = os.WriteFile(zedSettingsPath, []byte(updatedSettings), 0644)
	if err != nil {
		return fmt.Errorf("write updated Zed settings to %s: %w", zedSettingsPath, err)
	}

	log.Info("updated Zed theme", "theme", themeInfo.Name, "path", zedSettingsPath)

	return nil
}
