package match

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestMatchFuncByNameHit(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{
		"Generate": {{File: "g_test.go", TestFunc: "TestGen"}},
	}}
	fn := &model.Function{Name: "Generate"}
	refs, ok := MatchFuncByName(idx, nil, fn)
	if !ok || len(refs) != 1 || refs[0].TestFunc != "TestGen" {
		t.Fatalf("MatchFuncByName = %v ok=%v, want [TestGen]", refs, ok)
	}
}

func TestMatchFuncByNameMiss(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{"Generate": {{TestFunc: "TestGen"}}}}
	fn := &model.Function{Name: "Nope"}
	if refs, ok := MatchFuncByName(idx, nil, fn); ok || refs != nil {
		t.Errorf("MatchFuncByName(Nope) = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestMatchFuncByNameNilIndex(t *testing.T) {
	if refs, ok := MatchFuncByName(nil, nil, &model.Function{Name: "X"}); ok || refs != nil {
		t.Errorf("nil idx = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestMatchFuncByNameNilFunc(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{"X": {{TestFunc: "TestX"}}}}
	if refs, ok := MatchFuncByName(idx, nil, nil); ok || refs != nil {
		t.Errorf("nil fn = %v ok=%v, want nil,false", refs, ok)
	}
}

// bug003Funcs returns the matched TestFuncs for a method, running the full
// content-aware + source-receiver pipeline against a written package.
func bug003Funcs(t *testing.T, root, pkgDir string, fn *model.Function) []string {
	t.Helper()
	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	src, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	refs, ok := MatchFuncByName(idx, src, fn)
	if !ok {
		return nil
	}
	var out []string
	for _, r := range refs {
		out = append(out, r.TestFunc)
	}
	return out
}

// TestMatchFuncByNameReceiverAwareSeparation is the decisive BUG-003 case: two
// types implement a same-named method; each type's test instantiates only its
// own type via a composite literal. Each method must attribute only to its own
// test, never the other's.
func TestMatchFuncByNameReceiverAwareSeparation(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		"cs_file.go": `package gen
type CSharpFile struct{}
func (f *CSharpFile) GetFuncs() int { return 2 }
`,
		"go_file_get_funcs_test.go": `package gen
import "testing"
func TestGoFileGetFuncs(t *testing.T) { _ = (&GoFile{}).GetFuncs() }
`,
		"cs_file_get_funcs_test.go": `package gen
import "testing"
func TestCSharpFileGetFuncs(t *testing.T) {
	f := &CSharpFile{}
	_ = f.GetFuncs()
}
`,
	})

	goFn := &model.Function{Name: "GetFuncs", Receiver: "GoFile", File: pkgDir + "/go_file.go"}
	if got := bug003Funcs(t, root, pkgDir, goFn); len(got) != 1 || got[0] != "TestGoFileGetFuncs" {
		t.Fatalf("GoFile.GetFuncs -> %v, want [TestGoFileGetFuncs]", got)
	}
	csFn := &model.Function{Name: "GetFuncs", Receiver: "CSharpFile", File: pkgDir + "/cs_file.go"}
	if got := bug003Funcs(t, root, pkgDir, csFn); len(got) != 1 || got[0] != "TestCSharpFileGetFuncs" {
		t.Fatalf("CSharpFile.GetFuncs -> %v, want [TestCSharpFileGetFuncs] (composite via local var)", got)
	}
}

// TestMatchFuncByNameReceiverUntestedTypeUnmatched verifies that a type with no
// dedicated test for a same-name-multiple method is left unmatched (honest TODO)
// rather than mis-attributed to another type's test.
func TestMatchFuncByNameReceiverUntestedTypeUnmatched(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		"java_file.go": `package gen
type JavaFile struct{}
func (f *JavaFile) GetFuncs() int { return 2 }
`,
		// Only GoFile has a test.
		"go_file_get_funcs_test.go": `package gen
import "testing"
func TestGoFileGetFuncs(t *testing.T) { _ = (&GoFile{}).GetFuncs() }
`,
	})

	javaFn := &model.Function{Name: "GetFuncs", Receiver: "JavaFile", File: pkgDir + "/java_file.go"}
	if got := bug003Funcs(t, root, pkgDir, javaFn); got != nil {
		t.Fatalf("JavaFile.GetFuncs -> %v, want nil (no dedicated test, must not mis-attribute)", got)
	}
}

// TestMatchFuncByNameReceiverUnknownSingleRegression verifies the regression
// guard: when a method name has only one receiver in the package, a test whose
// call receiver cannot be statically resolved (constructor return) still
// attributes — the pre-existing correct behavior is preserved.
func TestMatchFuncByNameReceiverUnknownSingleRegression(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func NewGoFile() *GoFile { return &GoFile{} }
func (f *GoFile) GetFuncs() int { return 1 }
`,
		"go_file_get_funcs_test.go": `package gen
import "testing"
func TestGoFileGetFuncs(t *testing.T) { _ = NewGoFile().GetFuncs() }
`,
	})

	goFn := &model.Function{Name: "GetFuncs", Receiver: "GoFile", File: pkgDir + "/go_file.go"}
	if got := bug003Funcs(t, root, pkgDir, goFn); len(got) != 1 || got[0] != "TestGoFileGetFuncs" {
		t.Fatalf("GoFile.GetFuncs (unknown recv + single) -> %v, want [TestGoFileGetFuncs]", got)
	}
}

// TestMatchFuncByNameReceiverUnknownMultipleDropped verifies that an
// unresolvable receiver (constructor return) for a same-name-multiple method is
// dropped, so the method is not attributed to a test that may exercise another
// type.
func TestMatchFuncByNameReceiverUnknownMultipleDropped(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		"cs_file.go": `package gen
type CSharpFile struct{}
func NewCSharpFile() *CSharpFile { return &CSharpFile{} }
func (f *CSharpFile) GetFuncs() int { return 2 }
`,
		// This test calls GetFuncs via a constructor return -> receiver unknown.
		"ambiguous_test.go": `package gen
import "testing"
func TestAmbiguous(t *testing.T) { _ = NewCSharpFile().GetFuncs() }
`,
	})

	goFn := &model.Function{Name: "GetFuncs", Receiver: "GoFile", File: pkgDir + "/go_file.go"}
	if got := bug003Funcs(t, root, pkgDir, goFn); got != nil {
		t.Fatalf("GoFile.GetFuncs (unknown recv + multiple) -> %v, want nil (dropped)", got)
	}
}

// TestMatchFuncByNameReceiverHelperOneHop verifies helper 1-hop receiver
// consistency: the composite literal lives inside the helper, so the helper's
// own variable map resolves the receiver and attribution lands on the right
// type even though the test only calls the helper.
func TestMatchFuncByNameReceiverHelperOneHop(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		"cs_file.go": `package gen
type CSharpFile struct{}
func (f *CSharpFile) GetFuncs() int { return 2 }
`,
		"cs_via_helper_test.go": `package gen
import "testing"
func TestCSViaHelper(t *testing.T) { runCS(t) }
func runCS(t *testing.T) {
	f := &CSharpFile{}
	_ = f.GetFuncs()
}
`,
	})

	csFn := &model.Function{Name: "GetFuncs", Receiver: "CSharpFile", File: pkgDir + "/cs_file.go"}
	if got := bug003Funcs(t, root, pkgDir, csFn); len(got) != 1 || got[0] != "TestCSViaHelper" {
		t.Fatalf("CSharpFile.GetFuncs via helper -> %v, want [TestCSViaHelper]", got)
	}
	goFn := &model.Function{Name: "GetFuncs", Receiver: "GoFile", File: pkgDir + "/go_file.go"}
	if got := bug003Funcs(t, root, pkgDir, goFn); got != nil {
		t.Fatalf("GoFile.GetFuncs -> %v, want nil (helper instantiates only CSharpFile)", got)
	}
}

// TestMatchFuncByNameFreeFuncUnchanged verifies free functions still match by
// name regardless of receiver machinery (behavior unchanged / regression guard).
func TestMatchFuncByNameFreeFuncUnchanged(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"util.go": `package gen
func Parse(s string) int { return len(s) }
`,
		"parse_test.go": `package gen
import "testing"
func TestParse(t *testing.T) { _ = Parse("x") }
`,
	})

	fn := &model.Function{Name: "Parse", Receiver: "", File: pkgDir + "/util.go"}
	if got := bug003Funcs(t, root, pkgDir, fn); len(got) != 1 || got[0] != "TestParse" {
		t.Fatalf("free func Parse -> %v, want [TestParse]", got)
	}
}
