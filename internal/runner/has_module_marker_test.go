package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasModuleMarker_pom(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !hasModuleMarker(dir) {
		t.Error("expected true for dir with pom.xml")
	}
}

func TestHasModuleMarker_buildGradle(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "build.gradle"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if !hasModuleMarker(dir) {
		t.Error("expected true for dir with build.gradle")
	}
}

func TestHasModuleMarker_buildGradleKts(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "build.gradle.kts"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if !hasModuleMarker(dir) {
		t.Error("expected true for dir with build.gradle.kts")
	}
}

func TestHasModuleMarker_none(t *testing.T) {
	if hasModuleMarker(t.TempDir()) {
		t.Error("expected false for dir with no markers")
	}
}
