package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// installFakeJavaTool prepends a fake build-tool binary (with the given script
// body) to PATH and returns nothing; callers also create the matching marker.
func installFakeJavaTool(t *testing.T, name, body string) {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, name), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func mavenProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestJavaRunnerRunNoBuildTool(t *testing.T) {
	r := &JavaRunner{}
	_, err := r.Run(t.TempDir(), mkMatch("src/test/java/p/FooTest.java"))
	if err == nil {
		t.Fatal("expected error when no build tool marker is present")
	}
	if !strings.Contains(err.Error(), "build tool") {
		t.Errorf("error should mention build tool: %v", err)
	}
}

func TestJavaRunnerRunPass(t *testing.T) {
	dir := mavenProject(t)
	installFakeJavaTool(t, "mvn", "exit 0\n")

	r := &JavaRunner{}
	res, err := r.Run(dir, mkMatch("src/test/java/p/FooTest.java"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestJavaRunnerRunFail(t *testing.T) {
	dir := mavenProject(t)
	installFakeJavaTool(t, "mvn", "echo BUILD FAILURE 1>&2\nexit 1\n")

	r := &JavaRunner{}
	res, err := r.Run(dir, mkMatch("src/test/java/p/FooTest.java"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pass {
		t.Error("expected Pass=false when mvn test fails")
	}
	if res.Output == "" {
		t.Error("expected non-empty output for failing run")
	}
}

// TestJavaRunnerRunToolNotFound covers the findJavaTool error branch (line 25):
// a Maven project with no mvnw wrapper and no `mvn` on PATH.
func TestJavaRunnerRunToolNotFound(t *testing.T) {
	dir := mavenProject(t)
	// Point PATH at an empty dir so exec.LookPath("mvn") fails and there is no
	// mvnw wrapper in the project.
	t.Setenv("PATH", t.TempDir())

	r := &JavaRunner{}
	_, err := r.Run(dir, mkMatch("src/test/java/p/FooTest.java"))
	if err == nil {
		t.Fatal("expected error when mvn is not found on PATH")
	}
	if !strings.Contains(err.Error(), "not found on PATH") {
		t.Errorf("expected 'not found on PATH', got: %v", err)
	}
}

// TestJavaRunnerRunSubmoduleDir is the core BUG-004 fix: when the test file
// lives in a submodule with its own pom.xml (not in the root reactor), the
// runner must execute in the submodule directory, not the project root. The
// fake mvn records its working directory so we can assert on it.
func TestJavaRunnerRunSubmoduleDir(t *testing.T) {
	dir := mavenProject(t)
	sub := filepath.Join(dir, "examples")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	pwdFile := filepath.Join(t.TempDir(), "pwd")
	installFakeJavaTool(t, "mvn", "pwd > '"+pwdFile+"'\nexit 0\n")

	r := &JavaRunner{}
	if _, err := r.Run(dir, mkMatch("examples/src/test/java/p/FooTest.java")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(pwdFile)
	if err != nil {
		t.Fatalf("read pwd file: %v", err)
	}
	gotDir, err := filepath.EvalSymlinks(strings.TrimSpace(string(got)))
	if err != nil {
		t.Fatalf("eval pwd: %v", err)
	}
	wantDir, err := filepath.EvalSymlinks(sub)
	if err != nil {
		t.Fatalf("eval sub: %v", err)
	}
	if gotDir != wantDir {
		t.Errorf("mvn ran in %q, want submodule dir %q", gotDir, wantDir)
	}
}

func TestJavaRunnerImplementsRunner(t *testing.T) {
	var _ Runner = &JavaRunner{}
}
