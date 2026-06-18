package runner

import (
	"os"
	"path/filepath"
	"testing"
)

// writeMarker creates an empty build-marker file at projectRoot/relDir/name,
// making parent directories as needed. It is a test helper for laying out
// synthetic module trees.
func writeMarker(t *testing.T, root, relDir, name string) {
	t.Helper()
	dir := filepath.Join(root, relDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte("<project/>"), 0o644); err != nil {
		t.Fatalf("write marker %s: %v", name, err)
	}
}

// Single-module regression: a pom.xml only at the root resolves every file to
// projectRoot, preserving pre-Phase008 behavior.
func TestNearestModuleRoot_singleModuleRoot(t *testing.T) {
	root := t.TempDir()
	writeMarker(t, root, ".", "pom.xml")

	got := NearestModuleRoot(root, filepath.Join("src", "test", "java", "com", "example", "FooTest.java"))
	want := filepath.Clean(root)
	if got != want {
		t.Errorf("got %q, want projectRoot %q", got, want)
	}
}

// A submodule with its own pom.xml (not in the root <modules>) resolves files
// under it to the submodule root — the core BUG-004 fix.
func TestNearestModuleRoot_submodule(t *testing.T) {
	root := t.TempDir()
	writeMarker(t, root, ".", "pom.xml")
	writeMarker(t, root, "examples", "pom.xml")

	fileRel := filepath.Join("examples", "src", "test", "java", "com", "example", "FromJsonTest.java")
	got := NearestModuleRoot(root, fileRel)
	want := filepath.Join(root, "examples")
	if got != want {
		t.Errorf("got %q, want submodule root %q", got, want)
	}
}

// Nested poms: the innermost (nearest) marker wins.
func TestNearestModuleRoot_nestedNearestWins(t *testing.T) {
	root := t.TempDir()
	writeMarker(t, root, ".", "pom.xml")
	writeMarker(t, root, "examples", "pom.xml")
	writeMarker(t, filepath.Join(root, "examples"), "inner", "pom.xml")

	fileRel := filepath.Join("examples", "inner", "src", "main", "java", "Deep.java")
	got := NearestModuleRoot(root, fileRel)
	want := filepath.Join(root, "examples", "inner")
	if got != want {
		t.Errorf("got %q, want innermost module %q", got, want)
	}
}

// Gradle markers are recognized too.
func TestNearestModuleRoot_gradleSubmodule(t *testing.T) {
	root := t.TempDir()
	writeMarker(t, root, ".", "build.gradle")
	writeMarker(t, root, "lib", "build.gradle.kts")

	fileRel := filepath.Join("lib", "src", "test", "kotlin", "LibTest.kt")
	got := NearestModuleRoot(root, fileRel)
	want := filepath.Join(root, "lib")
	if got != want {
		t.Errorf("got %q, want gradle submodule root %q", got, want)
	}
}

// No marker found anywhere up to projectRoot falls back to projectRoot.
func TestNearestModuleRoot_noMarkerFallsBack(t *testing.T) {
	root := t.TempDir()
	// no markers at all
	got := NearestModuleRoot(root, filepath.Join("src", "Foo.java"))
	want := filepath.Clean(root)
	if got != want {
		t.Errorf("got %q, want projectRoot fallback %q", got, want)
	}
}

// A file directly at the project root (empty/dot dir) with a root marker
// resolves to projectRoot.
func TestNearestModuleRoot_fileAtRoot(t *testing.T) {
	root := t.TempDir()
	writeMarker(t, root, ".", "pom.xml")
	got := NearestModuleRoot(root, "Foo.java")
	want := filepath.Clean(root)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// A fileRel that climbs above projectRoot via ".." makes the walk start outside
// projectRoot, so it can never hit the root boundary; it must still terminate at
// the filesystem root and fall back to projectRoot rather than loop forever.
func TestNearestModuleRoot_escapesAboveRootFallsBack(t *testing.T) {
	root := t.TempDir() // no markers in any of its (temp) ancestors
	got := NearestModuleRoot(root, filepath.Join("..", "..", "..", "Foo.java"))
	want := filepath.Clean(root)
	if got != want {
		t.Errorf("got %q, want projectRoot fallback %q", got, want)
	}
}

// The upward walk stops at projectRoot even if a marker exists above it: a
// marker outside the project boundary is not selected.
func TestNearestModuleRoot_stopsAtProjectRoot(t *testing.T) {
	parent := t.TempDir()
	writeMarker(t, parent, ".", "pom.xml") // marker above the project
	root := filepath.Join(parent, "project")
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// project itself has no marker -> should fall back to project root, not parent.
	got := NearestModuleRoot(root, filepath.Join("src", "Foo.java"))
	want := filepath.Clean(root)
	if got != want {
		t.Errorf("got %q, want project root %q (must not escape to %q)", got, want, parent)
	}
}
