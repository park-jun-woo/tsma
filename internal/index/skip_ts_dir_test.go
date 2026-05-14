package index

import (
	"path/filepath"
	"testing"
)

func TestSkipTSDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "node_modules", path: "/project/node_modules", wantErr: true},
		{name: "dist", path: "/project/dist", wantErr: true},
		{name: "build", path: "/project/build", wantErr: true},
		{name: ".git", path: "/project/.git", wantErr: true},
		{name: ".tsma", path: "/project/.tsma", wantErr: true},
		{name: "src", path: "/project/src", wantErr: false},
		{name: "lib", path: "/project/lib", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := skipTSDir(tt.path)
			if tt.wantErr && err != filepath.SkipDir {
				t.Errorf("skipTSDir(%q) = %v, want filepath.SkipDir", tt.path, err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("skipTSDir(%q) = %v, want nil", tt.path, err)
			}
		})
	}
}
