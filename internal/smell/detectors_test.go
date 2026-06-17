//ff:func feature=smell type=test
//ff:what detect* 단위테스트: detectUnsafe/detectReflectDynamic/detectLinkname를 이름으로 직접 호출해 각 detector의 모든 분기(양성 매칭·음성 무발화)를 덮는다. parseDetectorSrc가 임시 _test.go를 go/ast로 파싱(주석 포함)해 *ast.File·fset·path를 돌려준다. ScanGo의 파싱-에러 경로도 덮는다.

package smell

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

// parseDetectorSrc parses src as a Go file (with comments) and returns the AST,
// its fset, and the file path so a detector can be exercised directly by name.
func parseDetectorSrc(t *testing.T, src string) (*ast.File, *token.FileSet, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture_test.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return file, fset, path
}

func TestDetectUnsafe_ImportAndSelector(t *testing.T) {
	src := `package p

import "unsafe"

func TestX() {
	var x int
	_ = unsafe.Pointer(&x)
}
`
	file, fset, path := parseDetectorSrc(t, src)
	got := detectUnsafe(file, fset, path)
	// One finding for the import, one for the unsafe.Pointer selector.
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2 (import + selector): %+v", len(got), got)
	}
	for _, f := range got {
		if f.Rule != "TS-REFL-001" {
			t.Errorf("Rule = %q, want TS-REFL-001", f.Rule)
		}
	}
}

func TestDetectUnsafe_Negative(t *testing.T) {
	// An import that is not unsafe and a selector whose base is not the unsafe
	// identifier: no findings, exercising the false sides of both guards.
	src := `package p

import "strings"

func TestX(s string) {
	_ = strings.TrimSpace(s)
}
`
	file, fset, path := parseDetectorSrc(t, src)
	if got := detectUnsafe(file, fset, path); len(got) != 0 {
		t.Fatalf("len = %d, want 0: %+v", len(got), got)
	}
}

func TestDetectReflectDynamic_MethodAndField(t *testing.T) {
	src := `package p

func TestX(v, obj R) {
	_ = v.MethodByName("m")
	_ = obj.FieldByName("f")
}
`
	file, fset, path := parseDetectorSrc(t, src)
	got := detectReflectDynamic(file, fset, path)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2 (MethodByName + FieldByName): %+v", len(got), got)
	}
	for _, f := range got {
		if f.Rule != "TS-REFL-002" {
			t.Errorf("Rule = %q, want TS-REFL-002", f.Rule)
		}
	}
}

func TestDetectReflectDynamic_Negative(t *testing.T) {
	// Legitimate reflect selectors must not fire (the not-in-set branch).
	src := `package p

import "reflect"

func TestX(a, b any) {
	_ = reflect.DeepEqual(a, b)
	_ = reflect.TypeOf(a)
}
`
	file, fset, path := parseDetectorSrc(t, src)
	if got := detectReflectDynamic(file, fset, path); len(got) != 0 {
		t.Fatalf("len = %d, want 0: %+v", len(got), got)
	}
}

func TestDetectLinkname_Directive(t *testing.T) {
	src := `package p

//go:linkname privateFn runtime.someInternal
func privateFn()
`
	file, fset, path := parseDetectorSrc(t, src)
	got := detectLinkname(file, fset, path)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1: %+v", len(got), got)
	}
	if got[0].Rule != "TS-REFL-003" {
		t.Errorf("Rule = %q, want TS-REFL-003", got[0].Rule)
	}
}

func TestDetectLinkname_Negative(t *testing.T) {
	// A prose comment that merely mentions go:linkname starts with "// " so the
	// directive prefix never matches (the false side of HasPrefix).
	src := `package p

// This comment mentions //go:linkname but is prose, not a directive.
func TestX() {}
`
	file, fset, path := parseDetectorSrc(t, src)
	if got := detectLinkname(file, fset, path); len(got) != 0 {
		t.Fatalf("len = %d, want 0: %+v", len(got), got)
	}
}

// TestScanGo_ParseError covers ScanGo's parse-error branch: a syntactically
// broken file yields (nil, err) so the caller can ignore it.
func TestScanGo_ParseError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken_test.go")
	if err := os.WriteFile(path, []byte("package p\nfunc ("), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := ScanGo(path)
	if err == nil {
		t.Fatal("expected a parse error for a broken file")
	}
	if got != nil {
		t.Errorf("findings = %+v, want nil on parse error", got)
	}
}
