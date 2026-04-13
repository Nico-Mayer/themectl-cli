package integrations

import (
	"fmt"
	"maps"
	"os"
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
	logger := log.New(os.Stderr)
	logger.SetPrefix(fmt.Sprintf("< %s >", i.Name()))
	return *logger
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
