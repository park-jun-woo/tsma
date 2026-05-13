package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Factory tests
// ---------------------------------------------------------------------------

func TestNewRunnerGo(t *testing.T) {
	r := NewRunner("go")
	if _, ok := r.(*GoRunner); !ok {
		t.Fatalf("NewRunner(\"go\") returned %T, want *GoRunner", r)
	}
}

func TestNewRunnerTS(t *testing.T) {
	r := NewRunner("typescript")
	if _, ok := r.(*TSRunner); !ok {
		t.Fatalf("NewRunner(\"typescript\") returned %T, want *TSRunner", r)
	}
}

func TestNewRunnerPython(t *testing.T) {
	r := NewRunner("python")
	if _, ok := r.(*PyRunner); !ok {
		t.Fatalf("NewRunner(\"python\") returned %T, want *PyRunner", r)
	}
}

func TestNewRunnerUnsupported(t *testing.T) {
	r := NewRunner("rust")
	u, ok := r.(*UnsupportedRunner)
	if !ok {
		t.Fatalf("NewRunner(\"rust\") returned %T, want *UnsupportedRunner", r)
	}
	if u.Lang != "rust" {
		t.Errorf("UnsupportedRunner.Lang = %q, want %q", u.Lang, "rust")
	}

	_, err := r.Run("/tmp", "main_test.rs")
	if err == nil {
		t.Fatal("UnsupportedRunner.Run should return an error")
	}
	if !strings.Contains(err.Error(), "rust") {
		t.Errorf("error should mention 'rust': %v", err)
	}
}

// ---------------------------------------------------------------------------
// extractTestFuncs tests
// ---------------------------------------------------------------------------

func TestExtractTestFuncs(t *testing.T) {
	dir := t.TempDir()
	content := `package foo

import "testing"

func TestLogin(t *testing.T) {
	// login test
}

func TestSignup(t *testing.T) {
	// signup test
}

func helperFunc() {
	// not a test
}

func BenchmarkLogin(b *testing.B) {
	// benchmark, not a Test func
}
`
	path := filepath.Join(dir, "handler_test.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := extractTestFuncs(path)
	if err != nil {
		t.Fatalf("extractTestFuncs: %v", err)
	}

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %v", len(funcs), funcs)
	}

	expected := map[string]bool{"TestLogin": true, "TestSignup": true}
	for _, f := range funcs {
		if !expected[f] {
			t.Errorf("unexpected func: %q", f)
		}
	}
}

func TestExtractTestFuncsEmpty(t *testing.T) {
	dir := t.TempDir()
	content := `package foo

func helperFunc() {}
`
	path := filepath.Join(dir, "helper.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := extractTestFuncs(path)
	if err != nil {
		t.Fatalf("extractTestFuncs: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected no funcs, got %v", funcs)
	}
}

func TestExtractTestFuncsNonexistent(t *testing.T) {
	_, err := extractTestFuncs("/nonexistent/path/test.go")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

// ---------------------------------------------------------------------------
// readFileBytes test
// ---------------------------------------------------------------------------

func TestReadFileBytes(t *testing.T) {
	dir := t.TempDir()
	content := "hello world"
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := readFileBytes(path)
	if err != nil {
		t.Fatalf("readFileBytes: %v", err)
	}
	if string(data) != content {
		t.Errorf("content = %q, want %q", string(data), content)
	}
}

func TestReadFileBytesNotFound(t *testing.T) {
	_, err := readFileBytes("/nonexistent/path/file.txt")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

// ---------------------------------------------------------------------------
// detectTSTestFramework tests
// ---------------------------------------------------------------------------

func writePkgJSON(t *testing.T, dir string, devDeps map[string]string) {
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

func TestDetectTSTestFrameworkVitest(t *testing.T) {
	dir := t.TempDir()
	writePkgJSON(t, dir, map[string]string{"vitest": "^1.0.0", "typescript": "^5.0.0"})

	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkJest(t *testing.T) {
	dir := t.TempDir()
	writePkgJSON(t, dir, map[string]string{"jest": "^29.0.0"})

	got := detectTSTestFramework(dir)
	if got != frameworkJest {
		t.Errorf("got %q, want %q", got, frameworkJest)
	}
}

func TestDetectTSTestFrameworkMocha(t *testing.T) {
	dir := t.TempDir()
	writePkgJSON(t, dir, map[string]string{"mocha": "^10.0.0"})

	got := detectTSTestFramework(dir)
	if got != frameworkMocha {
		t.Errorf("got %q, want %q", got, frameworkMocha)
	}
}

func TestDetectTSTestFrameworkFallback(t *testing.T) {
	dir := t.TempDir()
	writePkgJSON(t, dir, map[string]string{"typescript": "^5.0.0"})

	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (fallback)", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkNoPkgJSON(t *testing.T) {
	dir := t.TempDir()
	// No package.json at all — should fallback to vitest.
	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (fallback, no package.json)", got, frameworkVitest)
	}
}

func TestDetectTSTestFrameworkPriority(t *testing.T) {
	// When both vitest and jest are present, vitest wins (higher priority).
	dir := t.TempDir()
	writePkgJSON(t, dir, map[string]string{"vitest": "^1.0.0", "jest": "^29.0.0"})

	got := detectTSTestFramework(dir)
	if got != frameworkVitest {
		t.Errorf("got %q, want %q (vitest should have priority over jest)", got, frameworkVitest)
	}
}

// ---------------------------------------------------------------------------
// buildTestArgs tests
// ---------------------------------------------------------------------------

func TestBuildTestArgsVitest(t *testing.T) {
	args := buildTestArgs(frameworkVitest, "src/handler.test.ts")
	expected := []string{"vitest", "run", "src/handler.test.ts", "--reporter=verbose"}
	assertArgs(t, args, expected)
}

func TestBuildTestArgsJest(t *testing.T) {
	args := buildTestArgs(frameworkJest, "src/handler.test.ts")
	expected := []string{"jest", "src/handler.test.ts", "--verbose"}
	assertArgs(t, args, expected)
}

func TestBuildTestArgsMocha(t *testing.T) {
	args := buildTestArgs(frameworkMocha, "src/handler.test.ts")
	expected := []string{"mocha", "src/handler.test.ts"}
	assertArgs(t, args, expected)
}

func TestBuildTestArgsDefault(t *testing.T) {
	// Unknown framework should use vitest args.
	args := buildTestArgs("unknown", "src/handler.test.ts")
	expected := []string{"vitest", "run", "src/handler.test.ts", "--reporter=verbose"}
	assertArgs(t, args, expected)
}

func assertArgs(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %v (len %d), want %v (len %d)", got, len(got), want, len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("arg[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// detectPytest tests
// ---------------------------------------------------------------------------

func TestDetectPytestWithPyprojectToml(t *testing.T) {
	dir := t.TempDir()
	content := `[tool.pytest.ini_options]
minversion = "6.0"
`
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if !detectPytest(dir) {
		t.Error("expected true for pyproject.toml with [tool.pytest.ini_options]")
	}
}

func TestDetectPytestWithConftestPy(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "conftest.py"), []byte("# conftest"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !detectPytest(dir) {
		t.Error("expected true for conftest.py present")
	}
}

func TestDetectPytestWithPytestIni(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pytest.ini"), []byte("[pytest]"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !detectPytest(dir) {
		t.Error("expected true for pytest.ini present")
	}
}

func TestDetectPytestWithSetupCfg(t *testing.T) {
	dir := t.TempDir()
	content := `[tool:pytest]
addopts = -v
`
	if err := os.WriteFile(filepath.Join(dir, "setup.cfg"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if !detectPytest(dir) {
		t.Error("expected true for setup.cfg with [tool:pytest]")
	}
}

func TestDetectPytestWithRequirementsTxt(t *testing.T) {
	dir := t.TempDir()
	content := `flask==2.0.0
pytest==7.0.0
requests==2.28.0
`
	if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if !detectPytest(dir) {
		t.Error("expected true for requirements.txt with pytest")
	}
}

func TestDetectPytestWithRequirementsDevTxt(t *testing.T) {
	dir := t.TempDir()
	content := `pytest-cov==4.0.0
`
	if err := os.WriteFile(filepath.Join(dir, "requirements-dev.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if !detectPytest(dir) {
		t.Error("expected true for requirements-dev.txt with pytest")
	}
}

func TestDetectPytestNone(t *testing.T) {
	dir := t.TempDir()
	// No pytest indicators at all.
	if detectPytest(dir) {
		t.Error("expected false when no pytest indicators present")
	}
}

func TestDetectPytestEmptyPyprojectToml(t *testing.T) {
	dir := t.TempDir()
	content := `[tool.black]
line-length = 88
`
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// pyproject.toml exists but has no pytest section.
	if detectPytest(dir) {
		t.Error("expected false for pyproject.toml without pytest section")
	}
}
