package engine

import (
	"errors"
	"fmt"

	"github.com/nico-mayer/themectl-cli/internal/integration"
	"github.com/nico-mayer/themectl-cli/internal/theme"
)

type Engine struct {
	integrations []integration.Integration
}

func New(ints []integration.Integration) *Engine {
	return &Engine{
		integrations: ints,
	}
}

func (e *Engine) Apply(t theme.Resolved) error {
	var errs []error
	for _, in := range e.integrations {
		if err := in.Apply(t); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", in.Name(), err))
		}
	}
	return errors.Join(errs...)
}
