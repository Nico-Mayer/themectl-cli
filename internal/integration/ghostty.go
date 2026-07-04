package integration

import (
	"fmt"
	"os"
	"regexp"

	"github.com/nico-mayer/themectl-cli/internal/theme"
)

type Ghostty struct {
	ConfigPath string
}

var ghosttyThemeLine = regexp.MustCompile(`(?m)^(\s*theme\s*=\s*).*$`)

func setGhosttyTheme(config, themeName string) (string, error) {
	if !ghosttyThemeLine.MatchString(config) {
		return "", fmt.Errorf("no `theme =` setting found in ghostty config")
	}
	return ghosttyThemeLine.ReplaceAllString(config, `${1}"`+themeName+`"`), nil
}

func (Ghostty) Name() string {
	return "ghostty"
}

func (g Ghostty) Apply(t theme.Resolved) error {
	name, ok := t.Themes[g.Name()]
	if !ok {
		return fmt.Errorf("theme %s has no ghostty overwride", t.ID())
	}

	data, err := os.ReadFile(g.ConfigPath)
	if err != nil {
		return fmt.Errorf("read ghostty config: %w", err)
	}

	updated, err := setGhosttyTheme(string(data), name)
	if err != nil {
		return err
	}

	if err := os.WriteFile(g.ConfigPath, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write ghostty config: %w", err)
	}

	return nil
}
