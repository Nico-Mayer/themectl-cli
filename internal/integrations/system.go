package integrations

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

type SystemTheme struct{}

func (SystemTheme) Name() string {
	return "system-theme"
}

func (i SystemTheme) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	logger.Debug("applying system theme", "theme", themeInfo.Name, "appearance", themeInfo.Appearance, "os", runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		return setMacOSTheme(themeInfo, logger)
	default:
		logger.Warn("system theme integration is not supported on this operating system", "os", runtime.GOOS, "theme", themeInfo.Name)
		return fmt.Errorf("setting the system theme is not supported on %s", runtime.GOOS)
	}
}

func setMacOSTheme(themeInfo model.ThemeInfo, logger log.Logger) error {
	mode := strings.ToLower(themeInfo.Appearance)

	logger.Debug("preparing macOS appearance update", "theme", themeInfo.Name, "appearance", mode)

	var script string

	switch mode {
	case "dark":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to true`
	case "light":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to false`
	default:
		logger.Error("invalid appearance for macOS system theme integration", "theme", themeInfo.Name, "appearance", mode)
		return fmt.Errorf("unsupported appearance %q for theme %q: expected \"dark\" or \"light\"", mode, themeInfo.Name)
	}

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmedOutput := strings.TrimSpace(string(output))
		logger.Error("failed to apply macOS appearance", "theme", themeInfo.Name, "appearance", mode, "output", trimmedOutput, "err", err)
		return fmt.Errorf("failed to set macOS system appearance to %q for theme %q: %w (output: %s)", mode, themeInfo.Name, err, trimmedOutput)
	}

	logger.Info("applied macOS appearance", "theme", themeInfo.Name, "appearance", mode)
	return nil
}
