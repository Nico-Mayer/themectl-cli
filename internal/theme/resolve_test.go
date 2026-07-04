package theme

import (
	"reflect"
	"testing"
)

func TestResolve_VariantOverridesFamily(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}

	if got.Appearance != Light {
		t.Errorf("appearance = %q, want light", got.Appearance)
	}
	if !reflect.DeepEqual(got.WallpaperSources, []string{"catppuccin/macchiato"}) {
		t.Errorf("wallpaper sources = %v, want the variant's own list", got.WallpaperSources)
	}
	want := map[string]string{"ghostty": "catppuccin-latte", "eza": "cat-eza"}
	if !reflect.DeepEqual(got.Themes, want) {
		t.Errorf("themes = %v, want %v", got.Themes, want)
	}
	if got.ID() != "catppuccin/latte" {
		t.Errorf("id = %q", got.ID())
	}
}

func TestResolve_wallpaperSourcesNotInherited(t *testing.T) {
	fam := Family{Name: "f", Defaults: Spec{Appearance: new(Dark)}}
	v := Variant{Name: "v", Spec: Spec{}}
	got, _ := Resolve(fam, v)
	if len(got.WallpaperSources) != 0 {
		t.Errorf("want no sources, got %v", got.WallpaperSources)
	}
}

func TestResolve_missingAppearanceFails(t *testing.T) {
	fam := Family{Name: "f", Defaults: Spec{}}
	v := Variant{Name: "v", Spec: Spec{}}
	if _, err := Resolve(fam, v); err == nil {
		t.Fatal("expected error when appearance is set by neither family nor variant")
	}
}

func TestResolve_doesNotMutateInputs(t *testing.T) {
	fam := Family{Name: "f", Defaults: Spec{Appearance: new(Dark), Themes: map[string]string{"a": "1"}}}
	v := Variant{Name: "v", Spec: Spec{Themes: map[string]string{"b": "2"}}}
	_, _ = Resolve(fam, v)
	if len(fam.Defaults.Themes) != 1 {
		t.Errorf("Resolve mutated the family's map: %v", fam.Defaults.Themes)
	}
}
