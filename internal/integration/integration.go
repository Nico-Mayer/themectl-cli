package integration

import "github.com/Nico-Mayer/themectl/internal/theme"

type Integration interface {
	Name() string
	Apply(t theme.Resolved) error
}
