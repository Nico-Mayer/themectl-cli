package integration

import (
	"testing"

	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/testutil"
)

func TestEnabled_unknownNamesIgnored(t *testing.T) {
	cfg := config.Config{
		Settings: config.Settings{Integrations: []string{""}},
	}

	testutil.Equal(t, len(Enabled(cfg)), 0)
}
