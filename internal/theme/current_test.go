package theme

import (
	"path/filepath"
	"testing"
)

func TestCurrent_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".current")

	if _, err := ReadCurrent(path); err == nil {
		t.Error("expected error when no current theme is set")
	}

	want := "catppuccin/mocha"
	if err := WriteCurrent(path, want); err != nil {
		t.Fatal(err)
	}

	got, err := ReadCurrent(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
