package integrations

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/model"
)

type Integration interface {
	Name() string
	Apply(themeInfo model.ThemeInfo) error
}

func integrationLogger(i Integration) log.Logger {
	logger := log.New(os.Stderr)
	logger.SetPrefix(fmt.Sprintf("< %s >", i.Name()))
	return *logger
}
