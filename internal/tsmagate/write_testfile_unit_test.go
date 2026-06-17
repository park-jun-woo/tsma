//ff:func feature=gate type=test
//ff:what write_testfile 단위테스트: sanitizeGoSource가 마크다운 펜스(```go)·산문을 제거해 순수 Go만 남기는지, 펜스 없는 입력은 트림만 하는지, testTargetPath가 §2-2 우선순위(귀속 파일 → CanonicalTestPath 신규)대로 경로를 고르는지 검증한다.

package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

func TestSanitizeGoSource(t *testing.T) {
	cases := []struct {
		name, in string
		want     string // expected first line of the body
		reject   string // a substring that must NOT survive
	}{
		{
			name:   "fenced with go tag and prose",
			in:     "Here is your test:\n```go\npackage pkg\n\nfunc TestX(t *testing.T) {}\n```\nDone!",
			want:   "package pkg",
			reject: "Here is your test",
		},
		{
			name:   "bare fence no tag",
			in:     "```\npackage pkg\n```",
			want:   "package pkg",
			reject: "```",
		},
		{
			name:   "no fence just trims",
			in:     "\n\npackage pkg\nfunc TestX(t *testing.T) {}\n\n",
			want:   "package pkg",
			reject: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sanitizeGoSource(c.in)
			if !strings.HasPrefix(got, c.want) {
				t.Fatalf("got %q, want prefix %q", got, c.want)
			}
			if c.reject != "" && strings.Contains(got, c.reject) {
				t.Fatalf("output still contains rejected %q:\n%s", c.reject, got)
			}
			if !strings.HasSuffix(got, "\n") {
				t.Fatalf("output should end with a newline, got %q", got)
			}
		})
	}
}

func TestTestTargetPath_AttributedFileWins(t *testing.T) {
	p := funcPayload{Lang: "go", Root: "/r", Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	tm := match.TestMatch{Files: []string{filepath.Join("pkg", "foo_test.go")}}
	got, err := testTargetPath(p, tm, true)
	if err != nil {
		t.Fatalf("testTargetPath: %v", err)
	}
	if got != filepath.Join("pkg", "foo_test.go") {
		t.Fatalf("got %q, want the attributed file", got)
	}
}

func TestTestTargetPath_NewFuncDerivesCanonical(t *testing.T) {
	// No attributed test and no misnamed variant on disk (root does not exist):
	// falls through to CanonicalTestPath.
	p := funcPayload{Lang: "go", Root: t.TempDir(), Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	got, err := testTargetPath(p, match.TestMatch{}, false)
	if err != nil {
		t.Fatalf("testTargetPath: %v", err)
	}
	if want := filepath.Join("pkg", "foo_test.go"); got != want {
		t.Fatalf("got %q, want canonical %q", got, want)
	}
}

func TestTestTargetPath_AbsorbsMisnamed(t *testing.T) {
	// A misnamed test_foo_test.go exists but the canonical foo_test.go does not:
	// testTargetPath renames (absorbs) it and returns the canonical path.
	root := t.TempDir()
	pkgDir := filepath.Join(root, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "foo.go"), []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	misnamed := filepath.Join(pkgDir, "test_foo_test.go")
	if err := os.WriteFile(misnamed, []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	got, err := testTargetPath(p, match.TestMatch{}, false)
	if err != nil {
		t.Fatalf("testTargetPath: %v", err)
	}
	if want := filepath.Join("pkg", "foo_test.go"); got != want {
		t.Fatalf("got %q, want canonical %q after absorb", got, want)
	}
	// The misnamed file must have been renamed to the canonical name on disk.
	if _, err := os.Stat(misnamed); !os.IsNotExist(err) {
		t.Errorf("misnamed file should be gone after absorb, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(pkgDir, "foo_test.go")); err != nil {
		t.Errorf("canonical file should exist after absorb: %v", err)
	}
}

func TestTestTargetPath_RenameError(t *testing.T) {
	// The misnamed variant exists, but its directory is read-only so the absorb
	// rename fails and the error is surfaced (never silent).
	root := t.TempDir()
	pkgDir := filepath.Join(root, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "foo.go"), []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "test_foo_test.go"), []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(pkgDir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(pkgDir, 0o755) })

	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	if _, err := testTargetPath(p, match.TestMatch{}, false); err == nil {
		t.Fatal("expected a rename error when the directory is read-only")
	}
}

func TestTestTargetPath_UnderivablePathErrors(t *testing.T) {
	// A function in a language without a canonical-path rule, with no attributed
	// test and no misnamed variant: the path cannot be derived, so an error is
	// returned. (Python/TS/Rust are now handled, so use a still-unhandled language.)
	p := funcPayload{Lang: "ruby", Root: t.TempDir(), Fn: model.Function{File: filepath.Join("app", "x.rb")}}
	if _, err := testTargetPath(p, match.TestMatch{}, false); err == nil {
		t.Fatal("expected an error when no test path can be derived")
	}
}
