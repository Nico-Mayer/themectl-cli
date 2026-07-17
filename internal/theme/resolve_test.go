package theme

import (
	"testing"

	"github.com/nico-mayer/themectl-cli/internal/testutil"
)

func TestResolve_variantOverridesFamily(t *testing.T) {
	fam := Family{
		Name: "catppuccin",
		Defaults: Spec{
			Appearance: new(Dark),
			Themes:     map[string]string{"ghostty": "cat-default", "eza": "cat-eza"},
		},
	}
	v := Variant{
		Name: "latte",
		Spec: Spec{
			Appearance: new(Light),
			Themes:     map[string]string{"ghostty": "catppuccin-latte"},
		},
		WallpaperSources: []string{"catppuccin/macchiato"},
	}

	got, err := Resolve(fam, v)
	testutil.NoErr(t, err)
	testutil.Equal(t, got.Appearance, Light)
	testutil.Equal(t, got.ID(), "catppuccin/latte")
	testutil.Diff(t, []string{"catppuccin/macchiato"}, got.WallpaperSources)
	testutil.Diff(t, map[string]string{"ghostty": "catppuccin-latte", "eza": "cat-eza"}, got.Themes)
}

func TestResolve_variantInheritsAppearance(t *testing.T) {
	fam := Family{Name: "f", Defaults: Spec{Appearance: new(Dark)}}
	v := Variant{Name: "v"}

	got, err := Resolve(fam, v)
	testutil.NoErr(t, err)
	testutil.Equal(t, got.Appearance, Dark)
}

func TestResolve_wallpaperSourcesNotInherited(t *testing.T) {
	fam := Family{Name: "f", Defaults: Spec{Appearance: new(Dark)}}
	v := Variant{Name: "v"}

	got, err := Resolve(fam, v)
	testutil.NoErr(t, err)
	testutil.Equal(t, len(got.WallpaperSources), 0)
}

func TestResolve_missingAppearanceFails(t *testing.T) {
	if _, err := Resolve(Family{Name: "f"}, Variant{Name: "v"}); err == nil {
		t.Fatal("expected error when appearance is set by neither family nor variant")
	}
}

func TestResolve_doesNotMutateInputs(t *testing.T) {
	fam := Family{Name: "f", Defaults: Spec{Appearance: new(Dark), Themes: map[string]string{"a": "1"}}}
	v := Variant{Name: "v", Spec: Spec{Themes: map[string]string{"b": "2"}}}

	_, err := Resolve(fam, v)
	testutil.NoErr(t, err)
	testutil.Diff(t, map[string]string{"a": "1"}, fam.Defaults.Themes)
}
