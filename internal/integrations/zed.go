package integrations

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/nico-mayer/themectl-cli/internal/model"
)

type Zed struct{}

func (Zed) Name() string {
	return "zed"
}

func (i Zed) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	zedSettingsPath := filepath.Join(os.Getenv("HOME"), ".config", "zed", "settings.json")

	logger.Debug("updating Zed theme", "theme", themeInfo.Name, "path", zedSettingsPath)

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

	logger.Info("updated Zed theme", "theme", themeInfo.Name, "path", zedSettingsPath)

	return nil
}
