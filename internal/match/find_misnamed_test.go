package match

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeMatchFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFindMisnamedTestGoFound(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.go"))
	writeMatchFile(t, filepath.Join(srcDir, "test_foo_test.go"))

	misnamed, canonical, found := FindMisnamedTest(dir, "go", "pkg/foo.go")
	if !found {
		t.Fatal("expected found=true for test_foo_test.go variant")
	}
	if !strings.HasSuffix(canonical, "foo_test.go") || strings.HasSuffix(canonical, "test_foo_test.go") {
		t.Errorf("canonical = %q, want suffix foo_test.go", canonical)
	}
	if !strings.HasSuffix(misnamed, "test_foo_test.go") {
		t.Errorf("misnamed = %q, want suffix test_foo_test.go", misnamed)
	}
	wantMisnamed := filepath.Join("pkg", "test_foo_test.go")
	if misnamed != wantMisnamed {
		t.Errorf("misnamed = %q, want %q", misnamed, wantMisnamed)
	}
	wantCanonical := filepath.Join("pkg", "foo_test.go")
	if canonical != wantCanonical {
		t.Errorf("canonical = %q, want %q", canonical, wantCanonical)
	}
}

func TestFindMisnamedTestGoCanonicalExists(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.go"))
	writeMatchFile(t, filepath.Join(srcDir, "foo_test.go"))

	_, _, found := FindMisnamedTest(dir, "go", "pkg/foo.go")
	if found {
		t.Error("expected found=false when canonical foo_test.go exists")
	}
}

func TestFindMisnamedTestGoProductionSourceExcluded(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.go"))
	writeMatchFile(t, filepath.Join(srcDir, "test_foo.go"))

	_, _, found := FindMisnamedTest(dir, "go", "pkg/foo.go")
	if found {
		t.Error("expected found=false for test_foo.go (production source, not a test variant)")
	}
}

func TestFindMisnamedTestGoNoVariant(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.go"))

	_, _, found := FindMisnamedTest(dir, "go", "pkg/foo.go")
	if found {
		t.Error("expected found=false when no test variant exists")
	}
}

func TestFindMisnamedTestPythonExcluded(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.py"))
	writeMatchFile(t, filepath.Join(srcDir, "test_foo.py"))

	_, _, found := FindMisnamedTest(dir, "python", "pkg/foo.py")
	if found {
		t.Error("expected found=false for python (test_ prefix is canonical)")
	}
}

func TestFindMisnamedTestGoNonGoSourceFile(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.py"))
	writeMatchFile(t, filepath.Join(srcDir, "test_foo_test.go"))

	for _, sourceFile := range []string{"pkg/foo.py", "pkg/foo.txt"} {
		misnamed, canonical, found := FindMisnamedTest(dir, "go", sourceFile)
		if found {
			t.Errorf("sourceFile %q: expected found=false for non-.go source", sourceFile)
		}
		if misnamed != "" {
			t.Errorf("sourceFile %q: misnamed = %q, want empty", sourceFile, misnamed)
		}
		if canonical != "" {
			t.Errorf("sourceFile %q: canonical = %q, want empty", sourceFile, canonical)
		}
	}
}

func TestFindMisnamedTestOtherLangExcluded(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	writeMatchFile(t, filepath.Join(srcDir, "foo.go"))
	writeMatchFile(t, filepath.Join(srcDir, "test_foo_test.go"))

	for _, lang := range []string{"typescript", "rust"} {
		if _, _, found := FindMisnamedTest(dir, lang, "pkg/foo.go"); found {
			t.Errorf("lang %q: expected found=false", lang)
		}
	}
}
