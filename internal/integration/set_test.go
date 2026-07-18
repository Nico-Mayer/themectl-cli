package integration

import (
	"testing"

	"github.com/Nico-Mayer/themectl/internal/config"
	"github.com/Nico-Mayer/themectl/internal/testutil"
)

func TestEnabled_unknownNamesIgnored(t *testing.T) {
	cfg := config.Config{
		Settings: config.Settings{Integrations: []string{""}},
	}

	testutil.Equal(t, len(Enabled(cfg)), 0)
}
