package receipt

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
)

func Path(integrationName, themeName string) string {
	cfg, _ := config.Get()
	return filepath.Join(cfg.ReceiptsDir, integrationName, themeName)
}

func Load[T any](integrationName, themeName string) (T, error) {
	var data T
	target := filepath.Join(Path(integrationName, themeName), "receipt.json")
	bytes, err := os.ReadFile(target)
	if err != nil {
		return data, err
	}
	return data, json.Unmarshal(bytes, &data)
}

func Save[T any](integrationName, themeName string, data T) error {
	targetDir := Path(integrationName, themeName)

	err := os.MkdirAll(targetDir, 0o755)
	if err != nil {
		return err
	}

	byteData, err := json.Marshal(data)

	log.Debug("write", "receipt", string(byteData))
	err = os.WriteFile(filepath.Join(targetDir, "receipt.json"), byteData, 0o755)
	if err != nil {
		return err
	}

	return nil
}
