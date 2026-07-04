package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "themectl.toml"), []byte(`default-theme = "catppuccin-mocha"`), 0o744)

	c, err := Load(dir)
	if err != nil {
		t.Errorf("loading config %v", err)
	}

	wantTheme := "catppuccin-mocha"
	got := c.Settings.DefaultTheme
	if wantTheme != got {
		t.Errorf("want: %v got:%v", wantTheme, got)
	}
}
