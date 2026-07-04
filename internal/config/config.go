package config

import "path/filepath"

type Config struct {
	Root     string
	Settings Settings
}

func Load(root string) (Config, error) {
	s, err := loadSettings(filepath.Join(root, "themectl.toml"))
	if err != nil {
		return Config{}, err
	}

	return Config{Root: root, Settings: s}, nil
}
