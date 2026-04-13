package integrations

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/nico-mayer/themectl-cli/internal/model"
)

type Ghostty struct{}

func (Ghostty) Name() string {
	return "ghostty"
}

func (i Ghostty) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	ghosttyConfigPath := filepath.Join(os.Getenv("HOME"), ".config", "ghostty", "config.ghostty")

	if len(themeInfo.GhosttyThemeName) == 0 {
		themeInfo.GhosttyThemeName = themeInfo.Name
	}

	logger.Debug("updating theme", "ghostty_theme", themeInfo.GhosttyThemeName)

	data, err := os.ReadFile(ghosttyConfigPath)
	if err != nil {
		return fmt.Errorf("read Ghostty config %q: %w", ghosttyConfigPath, err)
	}
	content := string(data)

	re := regexp.MustCompile(`(theme\s*=\s*")([^"]*)(")`)
	if !re.MatchString(content) {
		return fmt.Errorf("Ghostty config %q does not contain a theme setting", ghosttyConfigPath)
	}

	updatedSettings := re.ReplaceAllString(content, `${1}`+themeInfo.GhosttyThemeName+`${3}`)

	err = os.WriteFile(ghosttyConfigPath, []byte(updatedSettings), 0644)
	if err != nil {
		return fmt.Errorf("write Ghostty config %q: %w", ghosttyConfigPath, err)
	}

	logger.Info("theme applied", "ghostty_theme", themeInfo.GhosttyThemeName)

	cmd := exec.Command("pkill", "-SIGUSR2", "ghostty")
	if err := cmd.Run(); err != nil {
		logger.Warn("reload signal failed", "err", err)
	}

	return nil
}
