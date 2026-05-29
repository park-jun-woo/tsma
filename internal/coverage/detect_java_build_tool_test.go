package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectJavaBuildToolMaven(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := detectJavaBuildTool(dir); got != "maven" {
		t.Errorf("got %q, want maven", got)
	}
}

func TestDetectJavaBuildToolGradle(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "build.gradle"), []byte("plugins {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := detectJavaBuildTool(dir); got != "gradle" {
		t.Errorf("got %q, want gradle", got)
	}
}

func TestDetectJavaBuildToolGradleKts(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "build.gradle.kts"), []byte("plugins {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := detectJavaBuildTool(dir); got != "gradle" {
		t.Errorf("got %q, want gradle", got)
	}
}

func TestDetectJavaBuildToolNone(t *testing.T) {
	if got := detectJavaBuildTool(t.TempDir()); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
