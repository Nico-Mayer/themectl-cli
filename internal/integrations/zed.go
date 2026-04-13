package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
	"github.com/nico-mayer/themectl-cli/internal/config"
	"github.com/nico-mayer/themectl-cli/internal/model"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type Zed struct{}

type ZedThemeInfo struct {
	ExtensionUrl string `json:"extension-url"`
	Theme        string `json:"theme"`
}

type ExtensionManifest struct {
	ID          string   `toml:"id"`
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	SchemaVer   int      `toml:"schema_version"`
	Authors     []string `toml:"authors"`
	Description string   `toml:"description"`
	Repository  string   `toml:"repository"`
}

func init() {
	Register(Zed{})
}

func (Zed) Name() string {
	return "zed"
}

func (i Zed) Apply(themeInfo model.ThemeInfo) error {
	logger := integrationLogger(i)

	zedThemeInfo, err := loadZedThemeInfo()
	if err != nil {
		return err
	}

	err = cloneExtension(zedThemeInfo)
	if err != nil {
		return err
	}

	zedSettingsPath := filepath.Join(os.Getenv("HOME"), ".config", "zed", "settings.json")

	data, err := os.ReadFile(zedSettingsPath)
	if err != nil {
		return fmt.Errorf("read Zed settings from %s: %w", zedSettingsPath, err)
	}

	content := string(data)
	re := regexp.MustCompile(`("theme"\s*:\s*")([^"]*)(")`)
	if !re.MatchString(content) {
		return fmt.Errorf("update Zed theme in %s: could not find \"theme\" setting", zedSettingsPath)
	}

	updatedSettings := re.ReplaceAllString(content, `${1}`+zedThemeInfo.Theme+`${3}`)

	err = os.WriteFile(zedSettingsPath, []byte(updatedSettings), 0644)
	if err != nil {
		return fmt.Errorf("write updated Zed settings to %s: %w", zedSettingsPath, err)
	}

	logger.Info("theme applied")

	return nil
}

func loadZedThemeInfo() (ZedThemeInfo, error) {
	cfg, _ := config.Get()
	zedThemeFilePath := filepath.Join(cfg.Paths.CurrentThemeDir, "zed.json")

	data, err := os.ReadFile(zedThemeFilePath)
	if err != nil {
		return ZedThemeInfo{}, err
	}

	var themeInfo ZedThemeInfo
	err = json.Unmarshal(data, &themeInfo)
	if err != nil {
		return ZedThemeInfo{}, err
	}

	return themeInfo, err
}

func cloneExtension(info ZedThemeInfo) error {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:   fmt.Sprintf("https://%s", info.ExtensionUrl),
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("clone extension: %w", err)
	}

	wt, err := r.Worktree()
	if err != nil {
		return err
	}

	f, err := wt.Filesystem.Open("extension.toml")
	if err != nil {
		return fmt.Errorf("open extension.toml: %w", err)
	}
	defer f.Close()

	var manifest ExtensionManifest
	if _, err := toml.NewDecoder(f).Decode(&manifest); err != nil {
		return fmt.Errorf("parse extension.toml: %w", err)
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	targetDir := filepath.Join(userConfigDir, "Zed", "extensions", "installed", manifest.ID)

	if _, err := os.Stat(targetDir); err == nil {
		log.Info("already installed", "extension", manifest.ID)
		return nil
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create extension dir: %w", err)
	}

	fs := wt.Filesystem
	err = billyWalk(fs, "/", func(path string, isDir bool) error {
		dst := filepath.Join(targetDir, path)
		if isDir {
			return os.MkdirAll(dst, 0o755)
		}
		src, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, src)
		return err
	})

	return err
}

func billyWalk(fs billy.Filesystem, root string, fn func(path string, isDir bool) error) error {
	entries, err := fs.ReadDir(root)
	if err != nil {
		return err
	}
	for _, e := range entries {
		p := filepath.Join(root, e.Name())
		if err := fn(p, e.IsDir()); err != nil {
			return err
		}
		if e.IsDir() {
			if err := billyWalk(fs, p, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
