package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestTSRunnerRun_absError covers the filepath.Abs error branch (line 14):
// Abs fails only when os.Getwd() fails, forced by removing the cwd.
func TestTSRunnerRun_absError(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd: %v", err)
	}
	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	r := &TSRunner{}
	if _, err := r.Run("/proj", "rel.test.ts"); err == nil {
		t.Fatal("expected error when filepath.Abs fails")
	}
}

func TestTSRunnerDetectsFrameworkAndBuildsArgs(t *testing.T) {
	dir := t.TempDir()

	// Write package.json with jest
	pkg := struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}{DevDependencies: map[string]string{"jest": "^29.0.0"}}
	data, err := json.Marshal(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	framework := detectTSTestFramework(dir)
	if framework != frameworkJest {
		t.Fatalf("framework = %q, want %q", framework, frameworkJest)
	}

	args := buildTestArgs(framework, "src/handler.test.ts")
	if len(args) != 3 {
		t.Fatalf("args = %v, want len 3", args)
	}
	if args[0] != "jest" {
		t.Errorf("args[0] = %q, want \"jest\"", args[0])
	}
	if args[1] != "src/handler.test.ts" {
		t.Errorf("args[1] = %q, want \"src/handler.test.ts\"", args[1])
	}
	if args[2] != "--verbose" {
		t.Errorf("args[2] = %q, want \"--verbose\"", args[2])
	}
}

func installFakeNpx(t *testing.T, body string) {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "npx"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestTSRunnerRun_pass(t *testing.T) {
	dir := t.TempDir()
	installFakeNpx(t, "exit 0\n")

	r := &TSRunner{}
	res, err := r.Run(dir, "src/handler.test.ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestTSRunnerRun_fail(t *testing.T) {
	dir := t.TempDir()
	installFakeNpx(t, "echo failed 1>&2\nexit 1\n")

	r := &TSRunner{}
	res, err := r.Run(dir, "src/handler.test.ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pass {
		t.Error("expected Pass=false when npx test fails")
	}
	if res.Output == "" {
		t.Error("expected non-empty output for failing run")
	}
}

func TestTSRunnerRun_relProjectRootRelFallback(t *testing.T) {
	parent := t.TempDir()
	if err := os.MkdirAll(filepath.Join(parent, "proj"), 0o755); err != nil {
		t.Fatal(err)
	}
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	if err := os.Chdir(parent); err != nil {
		t.Fatal(err)
	}

	installFakeNpx(t, "exit 0\n")

	// Relative projectRoot vs the absolute test path makes filepath.Rel fail,
	// exercising the relTest = absTest fallback branch.
	r := &TSRunner{}
	res, err := r.Run("proj", "src/handler.test.ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestTSRunnerFallbackFramework(t *testing.T) {
	dir := t.TempDir()
	// No package.json -> vitest fallback
	framework := detectTSTestFramework(dir)
	if framework != frameworkVitest {
		t.Errorf("framework = %q, want %q (fallback)", framework, frameworkVitest)
	}

	args := buildTestArgs(framework, "test.ts")
	if args[0] != "vitest" {
		t.Errorf("args[0] = %q, want \"vitest\"", args[0])
	}
}
