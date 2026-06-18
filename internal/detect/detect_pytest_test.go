package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// --- existing-signal regression ---

func TestDetectPytestToolSection(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pyproject.toml"), "[tool.pytest.ini_options]\nminversion = \"6.0\"\n")
	if !DetectPytest(dir) {
		t.Error("expected true for [tool.pytest]")
	}
}

func TestDetectPytestSetupCfg(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "setup.cfg"), "[tool:pytest]\naddopts = -v\n")
	if !DetectPytest(dir) {
		t.Error("expected true for [tool:pytest] in setup.cfg")
	}
}

func TestDetectPytestRequirementsTxt(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements.txt"), "pytest==7.0\n")
	if !DetectPytest(dir) {
		t.Error("expected true for requirements.txt pytest")
	}
}

func TestDetectPytestRequirementsDevTxt(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements-dev.txt"), "pytest-cov==4.0\n")
	if !DetectPytest(dir) {
		t.Error("expected true for requirements-dev.txt pytest")
	}
}

func TestDetectPytestPytestIni(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pytest.ini"), "[pytest]\n")
	if !DetectPytest(dir) {
		t.Error("expected true for pytest.ini")
	}
}

func TestDetectPytestConftest(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "conftest.py"), "# conftest\n")
	if !DetectPytest(dir) {
		t.Error("expected true for conftest.py")
	}
}

// --- D1: PEP 621 declarations ---

func TestDetectPytestOptionalDependencies(t *testing.T) {
	dir := t.TempDir()
	content := `[project]
name = "rulecat"

[project.optional-dependencies]
dev = ["pytest>=8.0"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if !DetectPytest(dir) {
		t.Error("expected true for [project.optional-dependencies] pytest (BUG-001)")
	}
}

func TestDetectPytestDependencyGroups(t *testing.T) {
	dir := t.TempDir()
	content := `[project]
name = "rulecat"

[dependency-groups]
test = ["pytest>=8.0", "pytest-cov"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if !DetectPytest(dir) {
		t.Error("expected true for [dependency-groups] pytest")
	}
}

func TestDetectPytestProjectDependencies(t *testing.T) {
	dir := t.TempDir()
	content := `[project]
name = "rulecat"
dependencies = ["click", "pytest>=8.0"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if !DetectPytest(dir) {
		t.Error("expected true for [project] dependencies pytest")
	}
}

func TestDetectPytestProjectDependenciesMultiline(t *testing.T) {
	dir := t.TempDir()
	// multiline array under [project] dependencies; the array opener line names
	// "dependencies", and a later line names pytest while still in [project].
	content := `[project]
name = "rulecat"
dependencies = [
    "click",
    "pytest>=8.0",
]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if !DetectPytest(dir) {
		t.Error("expected true for multiline [project] dependencies pytest")
	}
}

// --- D1: tests/ layout ---

func TestDetectPytestTestsLayoutTestPrefix(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pyproject.toml"), "[project]\nname = \"x\"\n")
	writeFile(t, filepath.Join(dir, "tests", "test_thing.py"), "def test_x():\n    assert True\n")
	if !DetectPytest(dir) {
		t.Error("expected true for tests/test_*.py layout")
	}
}

func TestDetectPytestTestsLayoutTestSuffix(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tests", "thing_test.py"), "def test_x():\n    assert True\n")
	if !DetectPytest(dir) {
		t.Error("expected true for tests/*_test.py layout")
	}
}

func TestDetectPytestTestsLayoutNonTestFilesIgnored(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tests", "__init__.py"), "")
	writeFile(t, filepath.Join(dir, "tests", "helpers.py"), "x = 1\n")
	if DetectPytest(dir) {
		t.Error("expected false: tests/ has no test_*.py / *_test.py")
	}
}

// --- D1: venv probe ---

func TestDetectPytestVenvProbe(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".venv", "bin", "pytest"), "#!/bin/sh\n")
	if !DetectPytest(dir) {
		t.Error("expected true for .venv/bin/pytest")
	}
}

func TestDetectPytestVenvProbeNonDotVenv(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "venv", "bin", "pytest"), "#!/bin/sh\n")
	if !DetectPytest(dir) {
		t.Error("expected true for venv/bin/pytest")
	}
}

// --- regression: pure unittest / no pytest stays false ---

func TestDetectPytestPureUnittestFalse(t *testing.T) {
	dir := t.TempDir()
	content := `[project]
name = "rulecat"
dependencies = ["click"]

[project.optional-dependencies]
dev = ["black", "ruff"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	// a unittest-style test under tests/ that is NOT pytest-named
	writeFile(t, filepath.Join(dir, "tests", "runner.py"), "import unittest\n")
	if DetectPytest(dir) {
		t.Error("expected false for pure unittest project (no pytest signal)")
	}
}

func TestDetectPytestEmptyFalse(t *testing.T) {
	if DetectPytest(t.TempDir()) {
		t.Error("expected false for empty dir")
	}
}

func TestDetectPytestPyprojectWithoutPytestFalse(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pyproject.toml"), "[tool.black]\nline-length = 88\n")
	if DetectPytest(dir) {
		t.Error("expected false for pyproject.toml without pytest")
	}
}

// --- direct helper coverage ---

func TestContainsPytestDepMissingFile(t *testing.T) {
	if containsPytestDep(t.TempDir()) {
		t.Error("expected false when pyproject.toml is absent")
	}
}

func TestContainsPytestDepOptionalDepSubsection(t *testing.T) {
	dir := t.TempDir()
	content := `[project.optional-dependencies.test]
deps = ["pytest"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if !containsPytestDep(dir) {
		t.Error("expected true for [project.optional-dependencies.test] subsection pytest")
	}
}

func TestContainsPytestDepClosesBlockOnNewSection(t *testing.T) {
	dir := t.TempDir()
	// pytest appears in a non-dependency section only -> false
	content := `[project.optional-dependencies]
dev = ["black"]

[tool.ruff]
extend-select = ["pytest-style"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if containsPytestDep(dir) {
		t.Error("expected false: pytest token only in a non-dependency section")
	}
}

func TestContainsPytestDepArrayClosesThenPytestElsewhere(t *testing.T) {
	dir := t.TempDir()
	// a multiline [project] dependencies array (no pytest) that closes, then a
	// real pytest declaration in a later dependency table -> exercises the
	// array-closing branch and still returns true.
	content := `[project]
name = "x"
dependencies = [
    "click",
    "rich",
]

[project.optional-dependencies]
dev = ["pytest>=8.0"]
`
	writeFile(t, filepath.Join(dir, "pyproject.toml"), content)
	if !containsPytestDep(dir) {
		t.Error("expected true: pytest declared after a closed dependencies array")
	}
}

func TestHasPytestLayoutSkipsTestNamedDir(t *testing.T) {
	dir := t.TempDir()
	// a directory literally named test_x.py inside tests/ must be skipped, and
	// with no real test file present detection stays false.
	if err := os.MkdirAll(filepath.Join(dir, "tests", "test_x.py"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "tests", "notes.txt"), "x")
	if hasPytestLayout(dir) {
		t.Error("expected false: only a dir named test_x.py and a non-py file")
	}
}

func TestHasPytestLayoutNoTestsDir(t *testing.T) {
	if hasPytestLayout(t.TempDir()) {
		t.Error("expected false when tests/ absent")
	}
}

func TestHasPytestLayoutSkipsSubdirs(t *testing.T) {
	dir := t.TempDir()
	// a subdir named like a test file must not count (only files do)
	if err := os.MkdirAll(filepath.Join(dir, "tests", "test_pkg.py"), 0o755); err != nil {
		t.Fatal(err)
	}
	if hasPytestLayout(dir) {
		t.Error("expected false: tests/test_pkg.py is a directory, not a file")
	}
}

func TestProbePytestVenvAbsent(t *testing.T) {
	if probePytestVenv(t.TempDir()) {
		t.Error("expected false when no venv pytest present")
	}
}

func TestContainsPytestMissingFile(t *testing.T) {
	if containsPytest(filepath.Join(t.TempDir(), "nope.toml"), "pytest") {
		t.Error("expected false for missing file")
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "f")
	if fileExists(p) {
		t.Error("expected false before create")
	}
	writeFile(t, p, "x")
	if !fileExists(p) {
		t.Error("expected true after create")
	}
}
