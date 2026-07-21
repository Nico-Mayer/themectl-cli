package integration

import (
	"fmt"
	"os"

	"github.com/Nico-Mayer/themectl/internal/theme"
)

type VSCode struct {
	SettingsPath string
	Installer    ExtensionInstaller
}

func (v VSCode) Name() string { return "vscode" }

func (v VSCode) Apply(t theme.Resolved) error {
	spec := t.VSCode
	if spec == nil || spec.Theme == "" {
		return fmt.Errorf("theme %s has no vscode override", t.ID())
	}

	if v.Installer != nil {
		for _, url := range spec.Extensions {
			if err := v.Installer.Ensure(url); err != nil {
				return err
			}
		}
	}

	data, err := os.ReadFile(v.SettingsPath)
	if err != nil {
		return fmt.Errorf("read vscode settings: %w", err)
	}

	updated, err := setJSONCString(string(data), "workbench.colorTheme", spec.Theme)
	if err != nil {
		return err
	}

	if spec.IconTheme != "" {
		updated, err = setJSONCString(updated, "workbench.iconTheme", spec.IconTheme)
		if err != nil {
			return err
		}
	}

	if err := os.WriteFile(v.SettingsPath, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write vscode settings: %w", err)
	}

	return nil
}

func (v VSCode) Check() error {
	return checkConfigDir(v.Name(), v.SettingsPath)
}
