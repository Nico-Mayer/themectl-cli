package integration

import (
	"testing"

	"github.com/nico-mayer/themectl-cli/internal/config"
)

// TODO: more tests when additional integrations are implemented
func TestSet(t *testing.T) {
	cfg := config.Config{
		Settings: config.Settings{
			Integrations: []string{
				"",
			},
		},
	}

	integrations := Enabled(cfg)

	if len(integrations) != 0 {
		t.Errorf("got integration back, but none wanted")
	}
}
