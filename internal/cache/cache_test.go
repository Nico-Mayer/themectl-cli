package cache

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nico-Mayer/themectl/internal/testutil"
)

func TestPutGet(t *testing.T) {
	c := New(filepath.Join(t.TempDir(), "sub"))

	key := "https://github.com/x/y"
	testutil.NoErr(t, c.Put(key, []byte{'a', 'b', 'c', '1', '2', '3'}))

	got, ok := c.Get(key)
	testutil.Equal(t, bytes.Equal(got, []byte{'a', 'b', 'c', '1', '2', '3'}), true)
	testutil.Equal(t, ok, true)
}

func TestFresh(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	key := "k"

	testutil.Equal(t, c.Fresh(key, time.Hour), false)
	testutil.NoErr(t, c.Put(key, []byte([]byte("v"))))
	testutil.Equal(t, c.Fresh(key, time.Hour), true)
	backdate(t, dir, 2*time.Hour)
	testutil.Equal(t, c.Fresh(key, time.Hour), false)
}

func TestTouch(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	testutil.NoErr(t, c.Put("k", []byte("v")))
	backdate(t, dir, 2*time.Hour)

	testutil.NoErr(t, c.Touch("k"))
	testutil.Equal(t, c.Fresh("k", time.Hour), true)
	got, _ := c.Get("k")
	testutil.Equal(t, bytes.Equal([]byte("v"), got), true)
}

func TestClear(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	testutil.NoErr(t, c.Put("k", []byte("v")))
	testutil.NoErr(t, c.Clear())

	got, ok := c.Get("k")
	testutil.Equal(t, bytes.Equal([]byte(""), got), true)
	testutil.Equal(t, ok, false)
}

func TestImutable(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	testutil.NoErr(t, c.Put("k", []byte{'1', '2', '3'}))

	temp, _ := c.Get("k")
	temp[0] = 2

	got, _ := c.Get("k")
	testutil.Equal(t, bytes.Equal(got, []byte{'1', '2', '3'}), true)
}

func backdate(t *testing.T, dir string, age time.Duration) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	testutil.NoErr(t, err)
	testutil.Equal(t, len(files), 1)

	old := time.Now().Add(-age)
	testutil.NoErr(t, os.Chtimes(files[0], old, old))
}
