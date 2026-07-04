package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	t.Run("missing file returns defaults", func(t *testing.T) {
		s, err := loadSettings(filepath.Join(t.TempDir(), "nope.toml"))
		if err != nil {
			t.Fatal(err)
		}
		if len(s.Integrations) == 0 {
			t.Error("expected default integrations")
		}
	})

	t.Run("file values override defaults, rest kept", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "settings.toml")
		data := "default-theme = \"catppuccin\"\n\n[config-dirs]\nghostty = \"/custom\"\n"
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}

		s, err := loadSettings(path)
		if err != nil {
			t.Fatal(err)
		}
		if s.DefaultTheme != "catppuccin" {
			t.Errorf("DefaultTheme = %q, want %q", s.DefaultTheme, "catppuccin")
		}
		if got := s.ConfigDirs["ghostty"]; got != "/custom" {
			t.Errorf("ghostty dir = %q, want %q", got, "/custom")
		}
		if len(s.Integrations) == 0 {
			t.Error("default integrations lost during merge")
		}
	})

	t.Run("invalid toml errors", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "settings.toml")
		if err := os.WriteFile(path, []byte("= nope"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := loadSettings(path); err == nil {
			t.Error("expected parse error")
		}
	})
}

func TestDefaultConfigDirs(t *testing.T) {
	unix := defaultConfigDirs("/home/u", "")
	if got, want := unix["zed"], filepath.Join("/home/u", ".config", "zed"); got != want {
		t.Errorf("zed = %q, want %q", got, want)
	}

	win := defaultConfigDirs(`C:\Users\u`, `C:\Users\u\AppData\Roaming`)
	if got, want := win["yazi"], filepath.Join(`C:\Users\u\AppData\Roaming`, "yazi", "config"); got != want {
		t.Errorf("yazi = %q, want %q", got, want)
	}
	if got, want := win["ghostty"], filepath.Join(`C:\Users\u`, ".config", "ghostty"); got != want {
		t.Errorf("ghostty = %q, want %q", got, want)
	}
}

func TestConfigDirFor(t *testing.T) {
	t.Setenv("HOME", "/home/u")
	t.Setenv("THEMECTL_TEST_DIR", "/from-env")

	tests := []struct {
		name string
		dirs map[string]string
		key  string
		want string
	}{
		{name: "nil map", dirs: nil, key: "ghostty", want: ""},
		{name: "missing key", dirs: map[string]string{}, key: "ghostty", want: ""},
		{name: "blank value", dirs: map[string]string{"ghostty": "   "}, key: "ghostty", want: ""},
		{name: "plain path", dirs: map[string]string{"ghostty": "/etc/ghostty"}, key: "ghostty", want: "/etc/ghostty"},
		{name: "env var", dirs: map[string]string{"ghostty": "$THEMECTL_TEST_DIR"}, key: "ghostty", want: "/from-env"},
		{name: "bare tilde", dirs: map[string]string{"ghostty": "~"}, key: "ghostty", want: "/home/u"},
		{name: "tilde prefix", dirs: map[string]string{"ghostty": "~/x"}, key: "ghostty", want: "/home/u/x"},
		{name: "tilde mid-path untouched", dirs: map[string]string{"ghostty": "/a/~/b"}, key: "ghostty", want: "/a/~/b"},
	}

	for _, tt := range tests {
		s := Settings{ConfigDirs: tt.dirs}
		if got := s.ConfigDirFor(tt.key); got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}
