package detect

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
	if tool := detectJavaBuildTool(dir); tool != "maven" {
		t.Errorf("tool = %q, want maven", tool)
	}
}

func TestDetectJavaBuildToolGradle(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "build.gradle"), []byte("plugins {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if tool := detectJavaBuildTool(dir); tool != "gradle" {
		t.Errorf("tool = %q, want gradle", tool)
	}
}

func TestDetectJavaBuildToolGradleKts(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "build.gradle.kts"), []byte("plugins {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if tool := detectJavaBuildTool(dir); tool != "gradle" {
		t.Errorf("tool = %q, want gradle", tool)
	}
}

func TestDetectJavaBuildToolMavenPriority(t *testing.T) {
	// When both pom.xml and build.gradle exist, Maven wins.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "build.gradle"), []byte("plugins {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if tool := detectJavaBuildTool(dir); tool != "maven" {
		t.Errorf("tool = %q, want maven (pom.xml takes priority)", tool)
	}
}

func TestDetectJavaBuildToolNone(t *testing.T) {
	dir := t.TempDir()
	if tool := detectJavaBuildTool(dir); tool != "" {
		t.Errorf("tool = %q, want empty", tool)
	}
}
