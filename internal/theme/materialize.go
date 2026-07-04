package theme

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (s *Store) Materialize(themeId, targetDir string) error {
	family, variant, exists := strings.Cut(themeId, "/")
	if !exists {
		return fmt.Errorf("theme id %q: want \"family/variant\"", themeId)
	}

	assets, err := s.Assets(family, variant)
	if err != nil {
		return err
	}

	err = os.RemoveAll(targetDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(targetDir, 0o755)
	if err != nil {
		return err
	}

	for name, p := range assets {
		data, err := fs.ReadFile(s.fsys, p)
		if err != nil {
			return fmt.Errorf("read asset %s: %w", p, err)
		}

		err = os.WriteFile(filepath.Join(targetDir, name), data, 0o644)
		if err != nil {
			return fmt.Errorf("write asset %s: %w", name, err)
		}
	}

	return nil
}
