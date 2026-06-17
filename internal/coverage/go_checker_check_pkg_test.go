//ff:test feature=coverage
package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCheckPkg_MultiFuncOneRun(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module checkpkgtest\n\ngo 1.22\n")
	write("add.go", "package mathx\n\nfunc Add(a, b int) int {\n\treturn a + b\n}\n")
	write("sub.go", "package mathx\n\nfunc Sub(a, b int) int {\n\tif a > b {\n\t\treturn a - b\n\t}\n\treturn b - a\n}\n")
	write("add_test.go", "package mathx\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) { _ = Add(1, 2) }\n")
	write("sub_test.go", "package mathx\n\nimport \"testing\"\n\nfunc TestSub(t *testing.T) { _ = Sub(3, 1) }\n")

	// File paths are project-root-relative (as the indexer produces), matching
	// how the cover profile's normalized paths are compared in overlaps.
	add := &model.Function{Name: "Add", File: "add.go", StartLine: 3, EndLine: 5}
	sub := &model.Function{Name: "Sub", File: "sub.go", StartLine: 3, EndLine: 8}

	c := &GoChecker{}
	reports, err := c.CheckPkg(root, root, []*model.Function{add, sub}, []string{"TestAdd", "TestSub"})
	if err != nil {
		t.Fatalf("CheckPkg: %v", err)
	}
	if reports[add] == nil || reports[sub] == nil {
		t.Fatalf("missing reports: add=%v sub=%v", reports[add], reports[sub])
	}
	if !reports[add].AllCovered {
		t.Errorf("Add: want fully covered, got %.1f%%", reports[add].TotalPct)
	}
	if reports[sub].AllCovered {
		t.Errorf("Sub: want partial (not all covered), got AllCovered")
	}
}

func TestCheckPkg_EmptyFuncs(t *testing.T) {
	c := &GoChecker{}
	reports, err := c.CheckPkg(t.TempDir(), t.TempDir(), nil, nil)
	if err != nil {
		t.Fatalf("CheckPkg empty: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("want empty reports, got %d", len(reports))
	}
}

// TestCheckPkg_RelErrorReturnsError: when pkgDir cannot be made relative to
// projectRoot (one absolute, one relative), CheckPkg returns the relative-path
// resolution error before attempting any test run.
func TestCheckPkg_RelErrorReturnsError(t *testing.T) {
	fn := &model.Function{Name: "F", File: "f.go", StartLine: 1, EndLine: 1}
	c := &GoChecker{}
	// projectRoot relative, pkgDir absolute -> filepath.Rel errors.
	if _, err := c.CheckPkg("relative/root", "/abs/pkg", []*model.Function{fn}, nil); err == nil {
		t.Fatal("expected error when pkgDir cannot be made relative to projectRoot")
	}
}

func TestCheckPkg_CompileFailureReturnsError(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module checkpkgfail\n\ngo 1.22\n")
	write("x.go", "package badp\n\nfunc Broken() int { return 1 }\n")
	write("x_test.go", "package badp\n\nimport \"testing\"\n\nfunc TestBroken(t *testing.T) { _ = Broken(); undefinedSymbol() }\n")

	broken := &model.Function{Name: "Broken", File: "x.go", StartLine: 3, EndLine: 3}
	c := &GoChecker{}
	if _, err := c.CheckPkg(root, root, []*model.Function{broken}, []string{"TestBroken"}); err == nil {
		t.Fatal("expected error from compile-failing package")
	}
}

// TestCheckPkg_parseProfileError covers the parse-profile error branch (line
// 55): `go test` exits 0 but writes an unparseable coverprofile. A fake `go` on
// PATH emits a single >64KB line so parseCoverProfile's scanner errors.
func TestCheckPkg_parseProfileError(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module checkpkgparse\n\ngo 1.22\n")
	write("foo.go", "package mathx\n\nfunc Foo() int { return 1 }\n")
	write("foo_test.go", "package mathx\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) { _ = Foo() }\n")

	// Fake `go`: exit 0 but write a >64KB single line (no newline) to the path
	// given via -coverprofile=..., overflowing bufio.Scanner's token limit.
	binDir := t.TempDir()
	script := `#!/bin/sh
for a in "$@"; do
  case "$a" in
    -coverprofile=*) out="${a#-coverprofile=}" ;;
  esac
done
if [ -n "$out" ]; then
  awk 'BEGIN{ for(i=0;i<100000;i++) printf "a" }' > "$out"
fi
exit 0
`
	if err := os.WriteFile(filepath.Join(binDir, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	foo := &model.Function{Name: "Foo", File: "foo.go", StartLine: 3, EndLine: 3}
	c := &GoChecker{}
	_, err := c.CheckPkg(root, root, []*model.Function{foo}, []string{"TestFoo"})
	if err == nil {
		t.Fatal("expected error when the coverage profile is unparseable")
	}
	if !strings.Contains(err.Error(), "parse coverage profile") {
		t.Errorf("expected 'parse coverage profile', got: %v", err)
	}
}
