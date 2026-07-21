package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nico-Mayer/themectl/internal/testutil"
	"github.com/Nico-Mayer/themectl/internal/theme"
)

type fakeVSCodeInstaller struct {
	ids []string
	err error
}

func (f *fakeVSCodeInstaller) Ensure(id string) error {
	f.ids = append(f.ids, id)
	return f.err
}

func writeVSCodeSettings(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "settings.json")
	testutil.NoErr(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestVSCode_Apply(t *testing.T) {
	settings := writeVSCodeSettings(t, `{"workbench.colorTheme": "old"}`)

	installer := &fakeVSCodeInstaller{}
	v := VSCode{SettingsPath: settings, Installer: installer}
	res := theme.Resolved{
		Family:  "catppuccin",
		Variant: "mocha",
		VSCode: &theme.VSCodeSpec{
			Theme:      "Catppuccin Mocha",
			Extensions: []string{"catppuccin.catppuccin-vsc", "catppuccin.catppuccin-vsc-icons"},
		},
	}

	testutil.NoErr(t, v.Apply(res))

	out, _ := os.ReadFile(settings)
	if !strings.Contains(string(out), `"workbench.colorTheme": "Catppuccin Mocha"`) {
		t.Errorf("colorTheme not rewritten: %q", out)
	}
	testutil.Diff(t, []string{"catppuccin.catppuccin-vsc", "catppuccin.catppuccin-vsc-icons"}, installer.ids)
}

func TestVSCode_Apply_noOverrideFails(t *testing.T) {
	v := VSCode{SettingsPath: "unused"}
	if err := v.Apply(theme.Resolved{}); err == nil {
		t.Error("expected error when theme has no vscode override")
	}
}

func TestVSCode_Apply_iconTheme(t *testing.T) {
	settings := writeVSCodeSettings(t, `{"workbench.iconTheme": "material-icon-theme"}`)

	v := VSCode{SettingsPath: settings}
	res := theme.Resolved{VSCode: &theme.VSCodeSpec{Theme: "X", IconTheme: "X Icons"}}
	testutil.NoErr(t, v.Apply(res))

	out, _ := os.ReadFile(settings)
	if !strings.Contains(string(out), `"workbench.iconTheme": "X Icons"`) {
		t.Errorf("iconTheme not rewritten: %q", out)
	}
}

func TestVSCode_Apply_unsetIconThemeLeftAlone(t *testing.T) {
	settings := writeVSCodeSettings(t, `{"workbench.iconTheme": "material-icon-theme"}`)

	v := VSCode{SettingsPath: settings}
	testutil.NoErr(t, v.Apply(theme.Resolved{VSCode: &theme.VSCodeSpec{Theme: "X"}}))

	out, _ := os.ReadFile(settings)
	if !strings.Contains(string(out), `"workbench.iconTheme": "material-icon-theme"`) {
		t.Errorf("iconTheme should be untouched: %q", out)
	}
}

func TestVSCode_Apply_installerErrorAborts(t *testing.T) {
	settings := writeVSCodeSettings(t, `{"workbench.colorTheme": "old"}`)

	installer := &fakeVSCodeInstaller{err: os.ErrPermission}
	v := VSCode{SettingsPath: settings, Installer: installer}
	res := theme.Resolved{VSCode: &theme.VSCodeSpec{
		Theme:      "X",
		Extensions: []string{"catppuccin.catppuccin-vsc"},
	}}

	if err := v.Apply(res); err == nil {
		t.Fatal("expected installer error to propagate")
	}

	out, _ := os.ReadFile(settings)
	testutil.Equal(t, string(out), `{"workbench.colorTheme": "old"}`)
}

func TestVSCode_Apply_nilInstaller(t *testing.T) {
	settings := writeVSCodeSettings(t, `{"workbench.colorTheme": "old"}`)

	v := VSCode{SettingsPath: settings}
	res := theme.Resolved{VSCode: &theme.VSCodeSpec{
		Theme:      "X",
		Extensions: []string{"catppuccin.catppuccin-vsc"},
	}}
	testutil.NoErr(t, v.Apply(res))

	out, _ := os.ReadFile(settings)
	if !strings.Contains(string(out), `"workbench.colorTheme": "X"`) {
		t.Errorf("theme not rewritten: %q", out)
	}
}
