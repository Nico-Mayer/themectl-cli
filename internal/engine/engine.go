package engine

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Nico-Mayer/themectl/internal/integration"
	"github.com/Nico-Mayer/themectl/internal/theme"
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
	warnErrors := make([]error, len(e.integrations))
	errs := make([]error, len(e.integrations))
	var wg sync.WaitGroup
	for i, in := range e.integrations {
		wg.Go(func() {
			if hc, ok := in.(integration.HealthChecker); ok {
				if err := hc.Check(); err != nil {
					warnErrors[i] = err
					return
				}
			}
			slog.Debug("applying integration", "integration", in.Name())
			if err := in.Apply(t); err != nil {
				slog.Warn("integration failed", "integration", in.Name(), "err", err)
				errs[i] = fmt.Errorf("%s: %w", in.Name(), err)
			}
		})
	}
	wg.Wait()

	for i, err := range warnErrors {
		if err != nil {
			slog.Warn("integration unhealthy, skipping", "integration", e.integrations[i].Name(), "err", err)
		}
	}

	return errors.Join(errs...)
}
