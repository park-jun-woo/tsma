package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectPytestPyprojectToml(t *testing.T) {
	dir := t.TempDir()
	content := "[tool.pytest.ini_options]\nminversion = \"6.0\"\n"
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("expected true for pyproject.toml with [tool.pytest")
	}
}

func TestDetectPytestSetupCfg(t *testing.T) {
	dir := t.TempDir()
	content := "[tool:pytest]\naddopts = -v\n"
	if err := os.WriteFile(filepath.Join(dir, "setup.cfg"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("expected true for setup.cfg with [tool:pytest]")
	}
}

func TestDetectPytestRequirementsTxt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("pytest==7.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("expected true for requirements.txt with pytest")
	}
}

func TestDetectPytestRequirementsDevTxt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "requirements-dev.txt"), []byte("pytest-cov==4.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("expected true for requirements-dev.txt with pytest")
	}
}

func TestDetectPytestPytestIni(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pytest.ini"), []byte("[pytest]"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("expected true for pytest.ini")
	}
}

func TestDetectPytestConftestPy(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "conftest.py"), []byte("# conftest"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("expected true for conftest.py")
	}
}

func TestDetectPytestNoIndicators(t *testing.T) {
	dir := t.TempDir()
	if detectPytest(dir) {
		t.Error("expected false when no pytest indicators present")
	}
}

func TestDetectPytestPyprojectWithoutPytest(t *testing.T) {
	dir := t.TempDir()
	content := "[tool.black]\nline-length = 88\n"
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if detectPytest(dir) {
		t.Error("expected false for pyproject.toml without pytest section")
	}
}
