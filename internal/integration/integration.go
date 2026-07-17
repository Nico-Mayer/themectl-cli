package integration

import "github.com/nico-mayer/themectl-cli/internal/theme"

type Integration interface {
	Name() string
	Apply(t theme.Resolved) error
}
