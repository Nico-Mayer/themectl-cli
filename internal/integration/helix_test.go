package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nico-mayer/themectl-cli/internal/theme"
)

func TestSetHelixTheme_pure(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		theme   string
		want    string
		wantErr bool
	}{
		{name: "simple", in: `theme = "old"`, theme: "new", want: `theme = "new"`},
		{name: "loose spacing", in: `theme   =   "old"`, theme: "new", want: `theme   =   "new"`},
		{name: "preserves siblings",
			in:    "theme = \"old\"\n[editor]\nline-number = \"relative\"\n",
			theme: "catppuccin_mocha",
			want:  "theme = \"catppuccin_mocha\"\n[editor]\nline-number = \"relative\"\n"},
		{name: "unquoted value not matched", in: `theme = old`, theme: "new", wantErr: true},
		{name: "missing theme key", in: "[editor]\nline-number = \"relative\"\n", theme: "new", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setHelixTheme(tt.in, tt.theme)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHelix_Apply(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte(`theme = "old"`), 0o644); err != nil {
		t.Fatal(err)
	}

	h := Helix{ConfigPath: cfgPath}
	res := theme.Resolved{
		Family:  "catppuccin",
		Variant: "mocha",
		Themes:  map[string]string{"helix": "catppuccin_mocha"},
	}

	if err := h.Apply(res); err != nil {
		t.Fatal(err)
	}
	out, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(out), `theme = "catppuccin_mocha"`) {
		t.Errorf("config not rewritten: %q", out)
	}
}

func TestHelix_Apply_noOverride(t *testing.T) {
	h := Helix{ConfigPath: "unused"}
	if err := h.Apply(theme.Resolved{Themes: map[string]string{}}); err == nil {
		t.Error("expected error when theme has no helix override")
	}
}
