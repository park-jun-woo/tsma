package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writePackageJSON(t *testing.T, dir string, devDeps map[string]string) {
	t.Helper()
	pkg := struct {
		DevDependencies map[string]string `json:"devDependencies,omitempty"`
	}{DevDependencies: devDeps}
	data, err := json.Marshal(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectTSTestFrameworkDetectsVitest(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"vitest": "^1.0.0"})
	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkDetectsJest(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"jest": "^29.0.0"})
	got := detectTSTestFramework(dir)
	if got != frameworkJest {
		t.Errorf("got %q, want %q", got, frameworkJest)
	}
}

func TestDetectTSTestFrameworkDetectsMocha(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"mocha": "^10.0.0"})
	got := detectTSTestFramework(dir)
	if got != frameworkMocha {
		t.Errorf("got %q, want %q", got, frameworkMocha)
	}
}

func TestDetectTSTestFrameworkFallbackNoPackageJSON(t *testing.T) {
	dir := t.TempDir()
	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (fallback)", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkFallbackNoDeps(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"typescript": "^5.0.0"})
	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (fallback)", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkVitestPriority(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"vitest": "^1.0.0", "jest": "^29.0.0"})
	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (vitest has priority)", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (fallback for invalid JSON)", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkDependenciesField(t *testing.T) {
	dir := t.TempDir()
	// Write package.json with vitest in dependencies (not devDependencies)
	pkg := struct {
		Dependencies map[string]string `json:"dependencies"`
	}{Dependencies: map[string]string{"vitest": "^1.0.0"}}
	data, err := json.Marshal(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (should detect from dependencies too)", got, frameworkVitest)
	}
}
