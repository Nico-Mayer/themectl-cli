package integrations

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

func SetSystemTheme(themeInfo model.ThemeInfo) error {
	log.Debug("applying system theme", "theme", themeInfo.Name, "appearance", themeInfo.Appearance, "os", runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		return setMacOSTheme(themeInfo)
	default:
		log.Warn("system theme integration is not supported on this operating system", "os", runtime.GOOS, "theme", themeInfo.Name)
		return fmt.Errorf("setting the system theme is not supported on %s", runtime.GOOS)
	}
}

func setMacOSTheme(themeInfo model.ThemeInfo) error {
	mode := strings.ToLower(themeInfo.Appearance)

	log.Debug("preparing macOS appearance update", "theme", themeInfo.Name, "appearance", mode)

	var script string

	switch mode {
	case "dark":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to true`
	case "light":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to false`
	default:
		log.Error("invalid appearance for macOS system theme integration", "theme", themeInfo.Name, "appearance", mode)
		return fmt.Errorf("unsupported appearance %q for theme %q: expected \"dark\" or \"light\"", mode, themeInfo.Name)
	}

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmedOutput := strings.TrimSpace(string(output))
		log.Error("failed to apply macOS appearance", "theme", themeInfo.Name, "appearance", mode, "output", trimmedOutput, "err", err)
		return fmt.Errorf("failed to set macOS system appearance to %q for theme %q: %w (output: %s)", mode, themeInfo.Name, err, trimmedOutput)
	}

	log.Info("applied macOS appearance", "theme", themeInfo.Name, "appearance", mode)
	return nil
}
