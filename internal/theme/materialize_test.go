package theme

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestStore_Materialize(t *testing.T) {
	fsys := fstest.MapFS{
		"catppuccin/family.toml":        {Data: []byte("[defaults]\nappearance='dark'\n")},
		"catppuccin/zed.json":           {Data: []byte(`{"from":"family"}`)},
		"catppuccin/mocha/variant.toml": {Data: []byte("appearance='dark'\n")},
		"catppuccin/mocha/nvim.lua":     {Data: []byte("-- mocha")},
	}
	dest := filepath.Join(t.TempDir(), "current")

	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(dest, "stale.txt"), []byte("stale"), 0o644)

	store := NewStore(fsys)
	err := store.Materialize("catppuccin/mocha", dest)
	if err != nil {
		t.Fatal(err)
	}

	if b, _ := os.ReadFile(filepath.Join(dest, "zed.json")); string(b) != `{"from":"family"}` {
		t.Errorf("inherited zed.json not materialized: %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "nvim.lua")); string(b) != "-- mocha" {
		t.Errorf("variant nvim.lua not materialized: %q", b)
	}
	if _, err := os.Stat(filepath.Join(dest, "stale.txt")); err == nil {
		t.Error("stale file survived; dest must be rebuilt from scratch")
	}
}
