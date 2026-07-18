package theme

import (
	"cmp"
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"path"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

var reservedNames = []string{
	"family.toml",
	"variant.toml",
}

type Store struct {
	fsys fs.FS
}

func NewStore(fsys fs.FS) *Store {
	return &Store{
		fsys: fsys,
	}
}

func (s *Store) Resolve(id string) (Resolved, error) {
	famName, variantName, found := strings.Cut(id, "/")
	if !found {
		return Resolved{}, fmt.Errorf("theme id %q: want \"family/variant\"", id)
	}

	family, err := s.family(famName)
	if err != nil {
		return Resolved{}, err
	}

	variant, err := s.variant(famName, variantName)
	if err != nil {
		return Resolved{}, err
	}

	return Resolve(family, variant)
}

func (s *Store) IDs() ([]string, error) {
	families, err := s.allFamilies()
	if err != nil {
		return nil, err
	}

	var out []string
	for _, fam := range families {
		variants, err := s.listVariants(fam)
		if err != nil {
			return nil, err
		}
		for _, v := range variants {
			out = append(out, fam+"/"+v)
		}
	}
	return out, nil
}

func (s *Store) List(a Appearance) ([]Resolved, error) {
	all, err := s.resolveAll()
	if err != nil {
		return nil, err
	}

	var out []Resolved
	for _, r := range all {
		if r.Appearance == a || a == AnyAppearance {
			out = append(out, r)
		}
	}

	return out, nil
}

func (s *Store) PickRandom(a Appearance) (Resolved, error) {
	candidates, err := s.List(a)
	if err != nil {
		return Resolved{}, err
	}

	if len(candidates) == 0 {
		return Resolved{}, fmt.Errorf("no matching candidates found for appearance %v", a)
	}

	return candidates[rand.IntN(len(candidates))], nil
}

func (s *Store) listVariants(family string) ([]string, error) {
	entries, err := fs.ReadDir(s.fsys, family)
	if err != nil {
		return []string{}, fmt.Errorf("read family %q: %w", family, err)
	}

	var out []string
	for _, e := range entries {
		if e.IsDir() && !(e.Name() == "wallpaper") {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

func (s *Store) resolveAll() ([]Resolved, error) {
	families, err := s.allFamilies()
	if err != nil {
		return nil, err
	}

	var (
		mu  sync.Mutex
		out []Resolved
		wg  sync.WaitGroup
	)

	for _, name := range families {
		wg.Go(func() {
			resolved := s.resolveFamily(name)
			mu.Lock()
			out = append(out, resolved...)
			mu.Unlock()
		})
	}
	wg.Wait()

	slices.SortFunc(out, func(a, b Resolved) int {
		return cmp.Compare(a.ID(), b.ID())
	})
	return out, nil
}

func (s *Store) resolveFamily(name string) []Resolved {
	fam, err := s.family(name)
	if err != nil {
		slog.Debug("skipping unresolvable family", "family", name, "err", err)
		return nil
	}

	variants, err := s.listVariants(name)
	if err != nil {
		slog.Debug("skipping unreadable family", "family", name, "err", err)
		return nil
	}

	var out []Resolved
	for _, v := range variants {
		variant, err := s.variant(name, v)
		if err != nil {
			slog.Debug("skipping unresolvable theme", "theme", name+"/"+v, "err", err)
			continue
		}
		res, err := Resolve(fam, variant)
		if err != nil {
			slog.Debug("skipping unresolvable theme", "theme", name+"/"+v, "err", err)
			continue
		}
		out = append(out, res)
	}
	return out
}

func (s *Store) allFamilies() ([]string, error) {
	entries, err := fs.ReadDir(s.fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read themes root: %w", err)
	}

	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	slices.Sort(out)

	return out, nil
}

func (s *Store) family(name string) (Family, error) {
	var wrap FamilyFile

	err := s.decode(path.Join(name, "family.toml"), &wrap)
	if err != nil {
		return Family{}, err
	}

	fam := Family{
		Name:     name,
		Defaults: wrap.Defaults,
	}
	return fam, nil
}

func (s *Store) variant(family, name string) (Variant, error) {
	var v VariantFile

	err := s.decode(path.Join(family, name, "variant.toml"), &v)
	if err != nil {
		return Variant{}, err
	}

	return Variant{
		Name:             name,
		WallpaperSources: v.WallpaperSources,
		Spec:             v.Spec,
	}, nil
}

func (s *Store) decode(path string, v any) error {
	data, err := fs.ReadFile(s.fsys, path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err = toml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	return nil
}
