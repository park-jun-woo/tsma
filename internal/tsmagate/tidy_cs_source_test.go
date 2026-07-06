package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTidyCsSourceFormatterAbsent asserts that when the .NET SDK (dotnet) is not
// on PATH, tidyCsSource degrades to returning the input unchanged — best-effort,
// never required. The remaining branches are driven below with a fake dotnet
// script on PATH, so every branch is deterministically reachable without the SDK.
func TestTidyCsSourceFormatterAbsent(t *testing.T) {
	emptyPath(t)
	src := "namespace P;\npublic class FooTests\n{\n    public void T() { }\n}\n"
	if got := tidyCsSource(src); got != src {
		t.Errorf("without dotnet the source must be unchanged: %q", got)
	}
}

// TestTidyCsSourceFormatSuccess drives the full happy path with a fake dotnet
// that rewrites Generated.cs inside the --folder dir, asserting the formatted
// (read-back) contents are returned.
func TestTidyCsSourceFormatSuccess(t *testing.T) {
	// The folder is the last CLI arg; overwrite the file dotnet format would
	// have formatted in place.
	installFakeTool(t, "dotnet", "#!/bin/sh\nfor a in \"$@\"; do d=\"$a\"; done\nprintf 'FORMATTED' > \"$d/Generated.cs\"\n")
	if got := tidyCsSource("class X{}\n"); got != "FORMATTED" {
		t.Errorf("tidyCsSource = %q, want %q", got, "FORMATTED")
	}
}

// TestTidyCsSourceRunError covers the cmd.Run failure branch: a fake dotnet that
// exits non-zero degrades to the input unchanged.
func TestTidyCsSourceRunError(t *testing.T) {
	installFakeTool(t, "dotnet", "#!/bin/sh\nexit 3\n")
	src := "class X{}\n"
	if got := tidyCsSource(src); got != src {
		t.Errorf("run failure must return src unchanged, got %q", got)
	}
}

// TestTidyCsSourceEmptyResult covers the empty read-back branch: a fake dotnet
// that truncates Generated.cs degrades to the input unchanged.
func TestTidyCsSourceEmptyResult(t *testing.T) {
	installFakeTool(t, "dotnet", "#!/bin/sh\nfor a in \"$@\"; do d=\"$a\"; done\n: > \"$d/Generated.cs\"\n")
	src := "class X{}\n"
	if got := tidyCsSource(src); got != src {
		t.Errorf("empty formatter output must return src unchanged, got %q", got)
	}
}

// TestTidyCsSourceMkdirTempError covers the MkdirTemp failure branch by pointing
// TMPDIR at a nonexistent directory (dotnet is present, so LookPath passes).
func TestTidyCsSourceMkdirTempError(t *testing.T) {
	installFakeTool(t, "dotnet", "#!/bin/sh\nexit 0\n")
	t.Setenv("TMPDIR", filepath.Join(t.TempDir(), "missing", "sub"))
	src := "class X{}\n"
	if got := tidyCsSource(src); got != src {
		t.Errorf("MkdirTemp failure must return src unchanged, got %q", got)
	}
}

// makeLongTempDir builds (and creates) a nested directory whose absolute path is
// exactly target characters long, each component within NAME_MAX.
func makeLongTempDir(t *testing.T, target int) string {
	t.Helper()
	dir := t.TempDir()
	for len(dir) < target {
		rem := target - len(dir)
		n := rem - 1
		if n > 200 {
			n = 200
		}
		if rem-(n+1) == 1 {
			n-- // never leave exactly 1 char (a bare "/") for the next round
		}
		dir = filepath.Join(dir, strings.Repeat("a", n))
	}
	if len(dir) != target {
		t.Fatalf("built path length %d, want %d", len(dir), target)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestTidyCsSourceTempWriteError covers the WriteFile failure branch: TMPDIR is
// sized so MkdirTemp succeeds (tempdir path <= PATH_MAX-1) but appending
// "/Generated.cs" exceeds PATH_MAX, so the write fails with ENAMETOOLONG and the
// input is returned unchanged.
func TestTidyCsSourceTempWriteError(t *testing.T) {
	installFakeTool(t, "dotnet", "#!/bin/sh\nexit 0\n")
	// MkdirTemp appends "/tsma-cs-fmt-" (13 chars) plus a 1..10-digit random
	// suffix: dir length is 4084..4093 (< 4096), while dir+"/Generated.cs" adds
	// 13 more chars and always exceeds the limit.
	t.Setenv("TMPDIR", makeLongTempDir(t, 4070))
	src := "class X{}\n"
	if got := tidyCsSource(src); got != src {
		t.Errorf("temp write failure must return src unchanged, got %q", got)
	}
}
