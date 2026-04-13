package integrations

import (
	"fmt"
	"maps"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

var registry = map[string]Integration{}

type Integration interface {
	Name() string
	Apply(themeInfo model.ThemeInfo) error
}

func integrationLogger(i Integration) log.Logger {
	l := log.Default().WithPrefix(fmt.Sprintf("< %s >", i.Name()))
	return *l
}

func All() []Integration {
	return slices.Collect(maps.Values(registry))
}

func Register(i Integration) {
	registry[i.Name()] = i
}

func Get(name string) (Integration, bool) {
	i, ok := registry[name]
	return i, ok
}
