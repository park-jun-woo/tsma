package index

import (
	"path/filepath"
	"testing"
)

func TestSkipGoDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "vendor", path: "/project/vendor", wantErr: true},
		{name: ".git", path: "/project/.git", wantErr: true},
		{name: ".tsma", path: "/project/.tsma", wantErr: true},
		{name: "node_modules", path: "/project/node_modules", wantErr: true},
		{name: "src", path: "/project/src", wantErr: false},
		{name: "internal", path: "/project/internal", wantErr: false},
		{name: "pkg", path: "/project/pkg", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := skipGoDir(tt.path)
			if tt.wantErr && err != filepath.SkipDir {
				t.Errorf("skipGoDir(%q) = %v, want filepath.SkipDir", tt.path, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("skipGoDir(%q) = %v, want nil", tt.path, err)
			}
		})
	}
}
