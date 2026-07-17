package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func NoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func Diff(t *testing.T, want, got any) {
	t.Helper()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("mismatch (-want +got):\n%s", d)
	}
}
