package session

import (
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	got := Dir("/tmp/myproject")
	want := filepath.Join("/tmp/myproject", ".tsma")
	if got != want {
		t.Errorf("Dir = %q, want %q", got, want)
	}
}

func TestDirWithTrailingSlash(t *testing.T) {
	got := Dir("/tmp/myproject/")
	want := filepath.Join("/tmp/myproject/", ".tsma")
	if got != want {
		t.Errorf("Dir = %q, want %q", got, want)
	}
}

func TestDirRelativePath(t *testing.T) {
	got := Dir("myproject")
	want := filepath.Join("myproject", ".tsma")
	if got != want {
		t.Errorf("Dir = %q, want %q", got, want)
	}
}
