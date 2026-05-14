package index

import (
	"path/filepath"
	"testing"
)

func TestSkipPyDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "__pycache__", path: "/project/__pycache__", wantErr: true},
		{name: ".venv", path: "/project/.venv", wantErr: true},
		{name: "venv", path: "/project/venv", wantErr: true},
		{name: ".git", path: "/project/.git", wantErr: true},
		{name: ".tsma", path: "/project/.tsma", wantErr: true},
		{name: "src", path: "/project/src", wantErr: false},
		{name: "app", path: "/project/app", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := skipPyDir(tt.path)
			if tt.wantErr && err != filepath.SkipDir {
				t.Errorf("skipPyDir(%q) = %v, want filepath.SkipDir", tt.path, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("skipPyDir(%q) = %v, want nil", tt.path, err)
			}
		})
	}
}
