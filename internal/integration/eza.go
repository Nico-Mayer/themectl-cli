package integration

import (
	"fmt"

	"github.com/nico-mayer/themectl-cli/internal/theme"
)

type Eza struct{}

func (Eza) Name() string {
	return "eza"
}

func (i Eza) Apply(t theme.Resolved) error {
	fmt.Println("run eza")
	return nil
}
