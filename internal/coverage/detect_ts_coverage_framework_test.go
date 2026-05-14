package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectTSCoverageFrameworkDetectsVitest(t *testing.T) {
	dir := t.TempDir()
	pkg := map[string]interface{}{
		"devDependencies": map[string]string{"vitest": "^1.0.0"},
	}
	data, _ := json.Marshal(pkg)
	os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644)

	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q", got, coverVitest)
	}
}

func TestDetectTSCoverageFrameworkDetectsJest(t *testing.T) {
	dir := t.TempDir()
	pkg := map[string]interface{}{
		"devDependencies": map[string]string{"jest": "^29.0.0"},
	}
	data, _ := json.Marshal(pkg)
	os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644)

	got := detectTSCoverageFramework(dir)
	if got != coverJest {
		t.Errorf("got %q, want %q", got, coverJest)
	}
}

func TestDetectTSCoverageFrameworkNoPackageJSON(t *testing.T) {
	dir := t.TempDir()
	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q (fallback)", got, coverVitest)
	}
}

func TestDetectTSCoverageFrameworkInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte("{bad}"), 0o644)

	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q (fallback for invalid JSON)", got, coverVitest)
	}
}

func TestDetectTSCoverageFrameworkNeitherDep(t *testing.T) {
	dir := t.TempDir()
	pkg := map[string]interface{}{
		"devDependencies": map[string]string{"typescript": "^5.0.0"},
	}
	data, _ := json.Marshal(pkg)
	os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644)

	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q (default fallback)", got, coverVitest)
	}
}

func TestDetectTSCoverageFrameworkVitestInDependencies(t *testing.T) {
	dir := t.TempDir()
	pkg := map[string]interface{}{
		"dependencies": map[string]string{"vitest": "^1.0.0"},
	}
	data, _ := json.Marshal(pkg)
	os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644)

	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q (vitest in dependencies)", got, coverVitest)
	}
}
