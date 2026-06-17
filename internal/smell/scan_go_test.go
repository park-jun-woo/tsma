//ff:func feature=smell type=test
//ff:what ScanGo 단위테스트: §8 양성(unsafe import/셀렉터, MethodByName/FieldByName, //go:linkname 각 1 Finding)·음성(reflect.DeepEqual/TypeOf/ValueOf, 주석·문자열 리터럴 속 키워드 → 0)·§6 위양성0 회귀(tsma 자체 _test.go 트리 전체 Findings 0, 특히 reflect.DeepEqual 파일에서 TS-REFL-002가 0).

package smell

import (
	"os"
	"path/filepath"
	"testing"
)

// scanSrc writes src to a temp _test.go file and returns its findings.
func scanSrc(t *testing.T, src string) []Finding {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture_test.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	findings, err := ScanGo(path)
	if err != nil {
		t.Fatalf("ScanGo: %v", err)
	}
	return findings
}

// countRule counts findings with the given rule ID.
func countRule(findings []Finding, rule string) int {
	n := 0
	for _, f := range findings {
		if f.Rule == rule {
			n++
		}
	}
	return n
}

func TestScanGo_UnsafeImport(t *testing.T) {
	src := `package p
import "unsafe"
func TestX() { _ = 1 }
`
	got := scanSrc(t, src)
	if n := countRule(got, "TS-REFL-001"); n != 1 {
		t.Fatalf("TS-REFL-001 = %d, want 1 (findings: %+v)", n, got)
	}
}

func TestScanGo_UnsafeSelector(t *testing.T) {
	// Reference unsafe.Pointer; the parser does not require the import to exist.
	src := `package p
func TestX() { var x int; _ = unsafe.Pointer(&x) }
`
	got := scanSrc(t, src)
	if n := countRule(got, "TS-REFL-001"); n != 1 {
		t.Fatalf("TS-REFL-001 = %d, want 1 (findings: %+v)", n, got)
	}
}

func TestScanGo_MethodByName(t *testing.T) {
	src := `package p
func TestX(v R) { _ = v.MethodByName("x") }
`
	got := scanSrc(t, src)
	if n := countRule(got, "TS-REFL-002"); n != 1 {
		t.Fatalf("TS-REFL-002 = %d, want 1 (findings: %+v)", n, got)
	}
}

func TestScanGo_FieldByName(t *testing.T) {
	src := `package p
func TestX(obj R) { _ = obj.FieldByName("y") }
`
	got := scanSrc(t, src)
	if n := countRule(got, "TS-REFL-002"); n != 1 {
		t.Fatalf("TS-REFL-002 = %d, want 1 (findings: %+v)", n, got)
	}
}

func TestScanGo_Linkname(t *testing.T) {
	src := `package p

//go:linkname privateFn runtime.someInternal
func privateFn()
`
	got := scanSrc(t, src)
	if n := countRule(got, "TS-REFL-003"); n != 1 {
		t.Fatalf("TS-REFL-003 = %d, want 1 (findings: %+v)", n, got)
	}
}

// TestScanGo_NegativeNoFalsePositive is the core false-positive-zero guard:
// legitimate reflect idioms and keywords inside comments/string literals must
// produce no findings.
func TestScanGo_NegativeNoFalsePositive(t *testing.T) {
	src := `package p

import "reflect"

// This comment mentions unsafe and MethodByName and //go:linkname but is prose.
func TestX(a, b any) {
	_ = reflect.DeepEqual(a, b)
	_ = reflect.TypeOf(a)
	_ = reflect.ValueOf(a)
	s := "unsafe MethodByName FieldByName //go:linkname"
	_ = s
}
`
	got := scanSrc(t, src)
	if len(got) != 0 {
		t.Fatalf("expected 0 findings, got %d: %+v", len(got), got)
	}
}

// TestScanGo_SelfTreeFalsePositiveZero is the §6 regression: scanning every
// _test.go in tsma's own tree must yield zero findings. In particular the files
// using reflect.DeepEqual must not trip TS-REFL-002 (the most common false
// positive). It walks up from this package to the module root.
func TestScanGo_SelfTreeFalsePositiveZero(t *testing.T) {
	root := moduleRoot(t)
	var scanned int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "bak" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || len(path) < len("_test.go") ||
			path[len(path)-len("_test.go"):] != "_test.go" {
			return nil
		}
		scanned++
		findings, perr := ScanGo(path)
		if perr != nil {
			// A parse error is not a smell; the gate handles broken builds.
			return nil
		}
		if len(findings) != 0 {
			t.Errorf("self-tree false positive in %s: %+v", path, findings)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if scanned == 0 {
		t.Fatal("scanned 0 _test.go files; walk root is wrong")
	}
	t.Logf("scanned %d self-tree _test.go files, 0 findings", scanned)
}

// moduleRoot walks up from the cwd until it finds go.mod (the tsma module root).
func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking up from cwd")
		}
		dir = parent
	}
}
