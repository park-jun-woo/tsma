//ff:func feature=gate type=test
//ff:what loop import/AST 헬퍼 단위테스트: collectUsedPackages(한정자 식별자 수집·중첩 selector 무시), importIdent(별칭/blank/dot/경로base/unquote에러), filterImportSpecs(used·blank·dot 보존, 미사용 드롭), pruneUnusedImports(블록 솎기·빈블록 제거·비import 보존), parseTestFuncs(추출/파싱실패), goTestFuncName(테스트함수/메서드/비Test/비FuncDecl)를 디스크 없이 go/parser로 직접 덮는다.

package tsmagate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"testing"
)

// parseSrc parses src into an *ast.File, failing the test on a parse error.
func parseSrc(t *testing.T, src string) *ast.File {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return file
}

func TestCollectUsedPackages_QualifiersAndNestedSelectors(t *testing.T) {
	// fmt.Println -> X is the bare ident "fmt" (added). a.b.c parses as a nested
	// SelectorExpr: the outer's X is itself a SelectorExpr (skipped), the inner's X
	// is the bare ident "a" (added). So "fmt" and "a" are used; "b"/"c" are not.
	file := parseSrc(t, "package pkg\n\nimport \"fmt\"\n\nfunc F(a struct{ b struct{ c int } }) { fmt.Println(a.b.c) }\n")
	used := collectUsedPackages(file)
	if !used["fmt"] {
		t.Errorf("expected qualifier %q to be collected: %v", "fmt", used)
	}
	if !used["a"] {
		t.Errorf("expected inner selector base %q to be collected: %v", "a", used)
	}
	if used["b"] || used["c"] {
		t.Errorf("a nested selector's non-ident X must not be collected: %v", used)
	}
}

func TestImportIdent_AllForms(t *testing.T) {
	// Alias, blank, dot and an unaliased path are all parseable; the unquote-error
	// case needs a hand-built spec whose path is not a valid quoted string.
	file := parseSrc(t, "package pkg\n\nimport (\n\tf \"fmt\"\n\t_ \"embed\"\n\t. \"strings\"\n\t\"bytes\"\n)\n")
	want := map[string]string{"f \"fmt\"": "f", "_ \"embed\"": "_", ". \"strings\"": ".", "\"bytes\"": "bytes"}
	for _, decl := range file.Decls {
		gd := decl.(*ast.GenDecl)
		for _, spec := range gd.Specs {
			is := spec.(*ast.ImportSpec)
			key := is.Path.Value
			if is.Name != nil {
				key = is.Name.Name + " " + is.Path.Value
			}
			if got := importIdent(is); got != want[key] {
				t.Errorf("importIdent(%s) = %q, want %q", key, got, want[key])
			}
		}
	}
	// An unparseable path literal falls back to the raw Path.Value.
	bad := &ast.ImportSpec{Path: &ast.BasicLit{Value: "not-a-quoted-string"}}
	if got := importIdent(bad); got != "not-a-quoted-string" {
		t.Errorf("importIdent(bad path) = %q, want the raw value", got)
	}
}

func TestFilterImportSpecs_KeepsUsedBlankDotDropsRest(t *testing.T) {
	file := parseSrc(t, "package pkg\n\nimport (\n\t\"fmt\"\n\t\"strings\"\n\t_ \"embed\"\n\t. \"bytes\"\n)\n")
	specs := file.Decls[0].(*ast.GenDecl).Specs
	used := map[string]bool{"fmt": true} // strings is unused
	kept := filterImportSpecs(specs, used)
	var idents []string
	for _, s := range kept {
		idents = append(idents, importIdent(s.(*ast.ImportSpec)))
	}
	sort.Strings(idents)
	got := strings.Join(idents, ",")
	if want := ".,_,fmt"; got != want {
		t.Fatalf("kept idents = %q, want %q (used fmt + blank + dot, strings dropped)", got, want)
	}
}

func TestPruneUnusedImports_FiltersDropsEmptyKeepsOther(t *testing.T) {
	// Grouped block: fmt used, strings unused (filtered). Solo block: bytes unused
	// (block becomes empty -> dropped). A var decl and a func decl are non-import
	// and must survive untouched.
	src := "package pkg\n\nimport (\n\t\"fmt\"\n\t\"strings\"\n)\n\nimport \"bytes\"\n\nvar V = 1\n\nfunc F() { fmt.Println(V) }\n"
	file := parseSrc(t, src)
	used := map[string]bool{"fmt": true}
	pruneUnusedImports(file, used)

	var importPaths []string
	var sawVar, sawFunc bool
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				for _, s := range d.Specs {
					importPaths = append(importPaths, s.(*ast.ImportSpec).Path.Value)
				}
			}
			if d.Tok == token.VAR {
				sawVar = true
			}
		case *ast.FuncDecl:
			sawFunc = true
		}
	}
	if strings.Join(importPaths, ",") != "\"fmt\"" {
		t.Errorf("imports after prune = %v, want only \"fmt\"", importPaths)
	}
	if !sawVar || !sawFunc {
		t.Errorf("non-import decls must be preserved: sawVar=%v sawFunc=%v", sawVar, sawFunc)
	}
}

func TestParseTestFuncs_ExtractsTestFunctions(t *testing.T) {
	src := "package pkg\n\nimport \"testing\"\n\n" +
		"type S struct{}\n\n" +
		"func (s S) TestMethod() {}\n\n" +
		"func helper() {}\n\n" +
		"func TestAlpha(t *testing.T) {}\n\n" +
		"func TestBeta(t *testing.T) {}\n"
	funcs, ok := parseTestFuncs(src)
	if !ok {
		t.Fatal("a parseable source must return ok=true")
	}
	sort.Strings(funcs)
	if got := strings.Join(funcs, ","); got != "TestAlpha,TestBeta" {
		t.Fatalf("extracted funcs = %q, want TestAlpha,TestBeta (method/helper excluded)", got)
	}
}

func TestParseTestFuncs_RejectsUnparseable(t *testing.T) {
	funcs, ok := parseTestFuncs("package pkg\n\nfunc TestX( {")
	if ok || funcs != nil {
		t.Fatalf("truncated source must return (nil, false), got (%v, %v)", funcs, ok)
	}
}

func TestGoTestFuncName_ClassifiesDecls(t *testing.T) {
	// One file exercises every branch: a GenDecl (import — not a FuncDecl), a method
	// (Recv != nil), a non-Test func, and a real test function.
	file := parseSrc(t, "package pkg\n\nimport \"testing\"\n\n"+
		"type S struct{}\n\n"+
		"func (s S) TestMethod() {}\n\n"+
		"func helper() {}\n\n"+
		"func TestReal(t *testing.T) {}\n")
	results := map[string]bool{}
	var passed []string
	for _, decl := range file.Decls {
		name, ok := goTestFuncName(decl)
		if ok {
			passed = append(passed, name)
		}
		// Record by a stable label so we assert each branch was hit.
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			results[fn.Name.Name] = ok
		}
	}
	if len(passed) != 1 || passed[0] != "TestReal" {
		t.Fatalf("only TestReal must be accepted, got %v", passed)
	}
	if results["TestMethod"] {
		t.Error("a method (Recv != nil) must be rejected")
	}
	if results["helper"] {
		t.Error("a non-Test-prefixed func must be rejected")
	}
}
