package integration

import (
	"testing"

	"github.com/Nico-Mayer/themectl/internal/testutil"
)

func TestJSONCString(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		key     string
		value   string
		want    string
		wantErr bool
	}{
		{name: "simple", in: `{"theme": "old"}`, key: "theme", value: "One Dark", want: `{"theme": "One Dark"}`},
		{name: "loose spacing", in: `{"theme"   :   "old"}`, key: "theme", value: "new", want: `{"theme"   :   "new"}`},
		{name: "preserves siblings and comments",
			in:    "{\n  // pick a theme\n  \"theme\": \"old\",\n  \"vim_mode\": true\n}",
			key:   "theme",
			value: "Catppuccin Mocha",
			want:  "{\n  // pick a theme\n  \"theme\": \"Catppuccin Mocha\",\n  \"vim_mode\": true\n}"},
		{name: "icon_theme", in: `{"icon_theme": "old"}`, key: "icon_theme", value: "new", want: `{"icon_theme": "new"}`},
		{name: "missing key appended", in: `{"vim_mode": true}`, key: "theme", value: "new",
			want: "{\"vim_mode\": true,\n  \"theme\": \"new\"\n}"},
		{name: "missing key appended after trailing comma", in: "{\n  \"vim_mode\": true,\n}", key: "theme", value: "new",
			want: "{\n  \"vim_mode\": true,\n  \"theme\": \"new\"\n}"},
		{name: "missing key appended to empty object", in: `{}`, key: "theme", value: "new",
			want: "{\n  \"theme\": \"new\"\n}"},
		{name: "value with quotes escaped", in: `{}`, key: "theme", value: `say "hi"`,
			want: "{\n  \"theme\": \"say \\\"hi\\\"\"\n}"},
		{name: "replaced value with dollar sign", in: `{"theme": "old"}`, key: "theme", value: "a$1b",
			want: `{"theme": "a$1b"}`},
		{name: "leading comment with brace ignored", in: "// {settings}\n{\n  \"vim_mode\": true\n}", key: "theme", value: "new",
			want: "// {settings}\n{\n  \"vim_mode\": true,\n  \"theme\": \"new\"\n}"},
		{name: "no object", in: `// just a comment`, key: "theme", value: "new", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setJSONCString(tt.in, tt.key, tt.value)
			testutil.Equal(t, err != nil, tt.wantErr)
			if !tt.wantErr {
				testutil.Equal(t, got, tt.want)
			}
		})
	}
}
