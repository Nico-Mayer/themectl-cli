package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nico-mayer/themectl-cli/internal/theme"
)

func TestSetGhosttyTheme_pure(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		theme   string
		want    string
		wantErr bool
	}{
		{name: "Simple", in: "theme = old\nfont=mono\n", theme: "new", want: "theme = \"new\"\nfont=mono\n", wantErr: false},
		{name: "Quoted spacing", in: `theme     =    "old"`, theme: "new", want: `theme     =    "new"`, wantErr: false},
		{name: "missing", in: "font = mono\n", theme: "new", want: "", wantErr: true},
	}

	for _, tt := range tests {
		got, err := setGhosttyTheme(tt.in, tt.theme)
		if (err != nil) != tt.wantErr {
			t.Fatalf("%s: err = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestGhostty_Apply(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.ghostty")
	if err := os.WriteFile(cfgPath, []byte("theme = old"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := Ghostty{
		ConfigPath: cfgPath,
	}

	res := theme.Resolved{
		Family:  "catppuccin",
		Variant: "mocha",
		Themes: map[string]string{
			"ghostty": "catppuccin-mocha",
		},
	}

	if err := g.Apply(res); err != nil {
		t.Fatal(err)
	}
	out, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(out), `theme = "catppuccin-mocha"`) {
		t.Errorf("config not rewritten: %q", out)
	}
}
