package integrations

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

func ChangeGhosttyTheme(themeInfo model.ThemeInfo) error {
	ghosttyConfigPath := filepath.Join(os.Getenv("HOME"), ".config", "ghostty", "config.ghostty")

	if len(themeInfo.GhosttyThemeName) == 0 {
		themeInfo.GhosttyThemeName = themeInfo.Name
	}

	log.Debug("updating Ghostty theme",
		"theme", themeInfo.Name,
		"ghostty_theme", themeInfo.GhosttyThemeName,
		"config_path", ghosttyConfigPath,
	)

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

	log.Info("updated Ghostty theme",
		"theme", themeInfo.Name,
		"ghostty_theme", themeInfo.GhosttyThemeName,
	)

	cmd := exec.Command("pkill", "-SIGUSR2", "ghostty")
	if err := cmd.Run(); err != nil {
		log.Warn("failed to signal Ghostty reload",
			"theme", themeInfo.Name,
			"signal", "SIGUSR2",
			"error", err,
		)
	} else {
		log.Debug("signaled Ghostty to reload config",
			"theme", themeInfo.Name,
			"signal", "SIGUSR2",
		)
	}

	return nil
}
