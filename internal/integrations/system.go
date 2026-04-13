package integrations

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/nico-mayer/themectl-cli/internal/model"
)

type SystemTheme struct{}

func (SystemTheme) Name() string {
	return "system-theme"
}

func (i SystemTheme) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	logger.Debug("applying", "appearance", themeInfo.Appearance, "os", runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		return i.setMacOSTheme(themeInfo)
	default:
		return fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
}

func (i SystemTheme) setMacOSTheme(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)
	mode := strings.ToLower(themeInfo.Appearance)

	var script string
	switch mode {
	case "dark":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to true`
	case "light":
		script = `tell application "System Events" to tell appearance preferences to set dark mode to false`
	default:
		return fmt.Errorf("unsupported appearance %q: expected \"dark\" or \"light\"", mode)
	}

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript failed: %w (output: %s)", err, strings.TrimSpace(string(output)))
	}

	logger.Info("theme applied", "appearance", mode)
	return nil
}
