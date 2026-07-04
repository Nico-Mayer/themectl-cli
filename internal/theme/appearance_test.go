package theme

import "testing"

func TestAppearance(t *testing.T) {
	tests := []struct {
		key     string
		want    Appearance
		wantErr bool
	}{
		{"dark", Dark, false},
		{"DARK", Dark, false},
		{"   DarK  ", Dark, false},
		{"", "", true},
		{"blue", "", true},
	}

	for _, tt := range tests {
		got, err := ParseAppearance(tt.key)
		if got != tt.want {
			t.Errorf("ParseAppearance(%q) = %q, want %q", tt.key, got, tt.want)
		}
		if (err != nil) != tt.wantErr {
			t.Fatalf("ParseAppearance(%q) err = %v, wantErr %v", tt.key, err, tt.wantErr)
		}
	}
}
