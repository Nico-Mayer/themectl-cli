package integration

import (
	"testing"

	"github.com/Nico-Mayer/themectl/internal/cache"
	"github.com/Nico-Mayer/themectl/internal/testutil"
)

func TestCodeInstaller_installsAndCaches(t *testing.T) {
	var calls []string
	c := codeInstaller{
		cache:   cache.New(t.TempDir()),
		install: func(id string) error { calls = append(calls, id); return nil },
	}

	testutil.NoErr(t, c.Ensure("carppuccin.catppuccin.vsc"))
	testutil.NoErr(t, c.Ensure("carppuccin.catppuccin.vsc"))

	testutil.Diff(t, []string{"carppuccin.catppuccin.vsc"}, calls)
}
