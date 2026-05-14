package index

import "testing"

func TestPkgDirOf(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{name: "root file", file: "handler.go", want: ""},
		{name: "single dir", file: "internal/handler.go", want: "internal"},
		{name: "nested dir", file: "internal/api/handler.go", want: "internal/api"},
		{name: "deep nesting", file: "a/b/c/d.go", want: "a/b/c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pkgDirOf(tt.file)
			if got != tt.want {
				t.Errorf("pkgDirOf(%q) = %q, want %q", tt.file, got, tt.want)
			}
		})
	}
}
