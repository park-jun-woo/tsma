package index

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// xn builds a treesitter.Node with a field name, text, and children.
func xn(typ, field, text string, children ...*treesitter.Node) *treesitter.Node {
	return &treesitter.Node{Type: typ, Field: field, Text: text, SRow: 0, ERow: 1, Children: children}
}

func TestExtractTSFunctionsNilRoot(t *testing.T) {
	if got := extractTSFunctions(nil, "a.ts", "a"); got != nil {
		t.Errorf("extractTSFunctions(nil) = %+v, want nil", got)
	}
}

func TestTSFuncFromDecl(t *testing.T) {
	// named declaration -> ok.
	node := xn("function_declaration", "", "", xn("identifier", "name", "add"))
	fn, ok := tsFuncFromDecl(node, "svc", "svc/api.ts", true)
	if !ok {
		t.Fatal("tsFuncFromDecl: ok=false for named decl")
	}
	if fn.Name != "add" || fn.File != "svc/api.ts" || !fn.Exported || fn.Status != model.StatusTodo {
		t.Errorf("fn = %+v", fn)
	}
	if !strings.HasSuffix(fn.QualifiedName, "add") {
		t.Errorf("QualifiedName = %q", fn.QualifiedName)
	}

	// no name child -> not ok.
	if _, ok := tsFuncFromDecl(xn("function_declaration", "", ""), "svc", "f.ts", false); ok {
		t.Error("tsFuncFromDecl: ok=true for nameless decl")
	}
	// empty name text -> not ok.
	empty := xn("function_declaration", "", "", xn("identifier", "name", ""))
	if _, ok := tsFuncFromDecl(empty, "svc", "f.ts", false); ok {
		t.Error("tsFuncFromDecl: ok=true for empty name")
	}
}

func TestTSArrowFuncFromDeclarator(t *testing.T) {
	mk := func(valType, name string) *treesitter.Node {
		children := []*treesitter.Node{}
		if name != "" {
			children = append(children, xn("identifier", "name", name))
		} else {
			children = append(children, xn("identifier", "name", ""))
		}
		if valType != "" {
			children = append(children, xn(valType, "value", ""))
		}
		return xn("variable_declarator", "", "", children...)
	}

	// arrow function binding -> ok.
	fn, ok := tsArrowFuncFromDeclarator(mk("arrow_function", "classify"), "svc", "svc/a.ts", true)
	if !ok || fn.Name != "classify" || !fn.Exported {
		t.Fatalf("arrow declarator = %+v ok=%v", fn, ok)
	}
	// function_expression and function are also accepted.
	if _, ok := tsArrowFuncFromDeclarator(mk("function_expression", "f"), "svc", "a.ts", false); !ok {
		t.Error("function_expression not accepted")
	}
	if _, ok := tsArrowFuncFromDeclarator(mk("function", "f"), "svc", "a.ts", false); !ok {
		t.Error("function not accepted")
	}

	// not a variable_declarator -> false.
	if _, ok := tsArrowFuncFromDeclarator(xn("comment", "", ""), "svc", "a.ts", false); ok {
		t.Error("non-declarator accepted")
	}
	// no value field -> false.
	if _, ok := tsArrowFuncFromDeclarator(mk("", "f"), "svc", "a.ts", false); ok {
		t.Error("declarator without value accepted")
	}
	// value is not a function (a plain literal) -> false.
	if _, ok := tsArrowFuncFromDeclarator(mk("number", "f"), "svc", "a.ts", false); ok {
		t.Error("non-function value accepted")
	}
	// empty name -> false.
	if _, ok := tsArrowFuncFromDeclarator(mk("arrow_function", ""), "svc", "a.ts", false); ok {
		t.Error("empty-name declarator accepted")
	}
}

func TestCollectTSArrowFuncs(t *testing.T) {
	lexical := xn("lexical_declaration", "", "",
		xn("variable_declarator", "", "",
			xn("identifier", "name", "fn"),
			xn("arrow_function", "value", ""),
		),
		// a non-function declarator is skipped.
		xn("variable_declarator", "", "",
			xn("identifier", "name", "x"),
			xn("number", "value", ""),
		),
	)
	var out []model.Function
	collectTSArrowFuncs(lexical, "svc", "svc/a.ts", true, &out)
	if len(out) != 1 || out[0].Name != "fn" {
		t.Fatalf("collectTSArrowFuncs = %+v", out)
	}
}

func TestTSMethodFromDefinition(t *testing.T) {
	// exported (uppercase) method -> ok, Exported true.
	up := xn("method_definition", "", "", xn("property_identifier", "name", "Area"))
	fn, ok := tsMethodFromDefinition(up, "Rect", "svc", "svc/a.ts")
	if !ok || fn.Name != "Area" || !fn.Exported {
		t.Fatalf("uppercase method = %+v ok=%v", fn, ok)
	}
	if !strings.Contains(fn.QualifiedName, "Rect") {
		t.Errorf("QualifiedName lacks receiver: %q", fn.QualifiedName)
	}
	// lowercase method -> ok, Exported false.
	low := xn("method_definition", "", "", xn("property_identifier", "name", "scale"))
	fn, ok = tsMethodFromDefinition(low, "Rect", "svc", "a.ts")
	if !ok || fn.Exported {
		t.Errorf("lowercase method = %+v ok=%v", fn, ok)
	}
	// not a method_definition -> false.
	if _, ok := tsMethodFromDefinition(xn("public_field_definition", "", ""), "Rect", "svc", "a.ts"); ok {
		t.Error("non-method accepted")
	}
	// nameless -> false.
	if _, ok := tsMethodFromDefinition(xn("method_definition", "", ""), "Rect", "svc", "a.ts"); ok {
		t.Error("nameless method accepted")
	}
	// empty name -> false.
	emptyName := xn("method_definition", "", "", xn("property_identifier", "name", ""))
	if _, ok := tsMethodFromDefinition(emptyName, "Rect", "svc", "a.ts"); ok {
		t.Error("empty-name method accepted")
	}
	// constructor -> skipped.
	ctor := xn("method_definition", "", "", xn("property_identifier", "name", "constructor"))
	if _, ok := tsMethodFromDefinition(ctor, "Rect", "svc", "a.ts"); ok {
		t.Error("constructor accepted")
	}
}

func TestCollectTSMethods(t *testing.T) {
	class := xn("class_declaration", "", "",
		xn("identifier", "name", "Rect"),
		xn("class_body", "", "",
			xn("method_definition", "", "", xn("property_identifier", "name", "Area")),
			xn("method_definition", "", "", xn("property_identifier", "name", "constructor")),
		),
	)
	var out []model.Function
	collectTSMethods(class, "svc", "svc/a.ts", &out)
	if len(out) != 1 || out[0].Name != "Area" {
		t.Fatalf("collectTSMethods = %+v", out)
	}

	// a class with no class_body emits nothing (and does not panic).
	noBody := xn("class_declaration", "", "", xn("identifier", "name", "Empty"))
	var out2 []model.Function
	collectTSMethods(noBody, "svc", "a.ts", &out2)
	if len(out2) != 0 {
		t.Errorf("class without body emitted %+v", out2)
	}

	// an anonymous class (no name field) still works, className empty.
	anon := xn("class_declaration", "", "",
		xn("class_body", "", "",
			xn("method_definition", "", "", xn("property_identifier", "name", "Run")),
		),
	)
	var out3 []model.Function
	collectTSMethods(anon, "svc", "a.ts", &out3)
	if len(out3) != 1 || out3[0].Name != "Run" {
		t.Errorf("anonymous class = %+v", out3)
	}
}

func TestWalkTSChild(t *testing.T) {
	var out []model.Function
	relDir, relPath := "svc", "svc/a.ts"

	// export_statement recurses with exported=true.
	exp := xn("export_statement", "", "",
		xn("function_declaration", "", "", xn("identifier", "name", "Exported")),
	)
	walkTSChild(exp, relDir, relPath, false, &out)

	// plain function_declaration.
	walkTSChild(xn("function_declaration", "", "", xn("identifier", "name", "plain")), relDir, relPath, false, &out)

	// lexical_declaration with an arrow.
	walkTSChild(xn("lexical_declaration", "", "",
		xn("variable_declarator", "", "", xn("identifier", "name", "arrow"), xn("arrow_function", "value", "")),
	), relDir, relPath, false, &out)

	// class_declaration.
	walkTSChild(xn("class_declaration", "", "",
		xn("identifier", "name", "C"),
		xn("class_body", "", "", xn("method_definition", "", "", xn("property_identifier", "name", "M"))),
	), relDir, relPath, false, &out)

	// an unrecognized node type is ignored (default case).
	walkTSChild(xn("import_statement", "", ""), relDir, relPath, false, &out)

	got := map[string]bool{}
	for _, f := range out {
		got[f.Name] = f.Exported
	}
	for _, name := range []string{"Exported", "plain", "arrow", "M"} {
		if _, ok := got[name]; !ok {
			t.Errorf("missing %q in %v", name, got)
		}
	}
	if !got["Exported"] {
		t.Error("export_statement child not marked exported")
	}
	if got["plain"] {
		t.Error("plain function wrongly marked exported")
	}
}

func TestWalkTSTopLevel(t *testing.T) {
	program := xn("program", "", "",
		xn("function_declaration", "", "", xn("identifier", "name", "a")),
		xn("function_declaration", "", "", xn("identifier", "name", "b")),
	)
	var out []model.Function
	walkTSTopLevel(program, "svc", "svc/a.ts", false, &out)
	if len(out) != 2 {
		t.Fatalf("walkTSTopLevel = %+v", out)
	}
}

func TestCollectSourceFiles(t *testing.T) {
	root := t.TempDir()
	mustWrite := func(rel, body string) {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite(".tsmignore", "ignored/\nsecret.src\n")
	mustWrite("a.src", "x")          // indexable source
	mustWrite("b.txt", "x")          // not a source
	mustWrite("secret.src", "x")     // an ignored *file* -> skipped (not a dir)
	mustWrite("sub/d.src", "x")      // indexable source in a walked dir
	mustWrite("ignored/c.src", "x")  // under an ignored dir -> SkipDir
	mustWrite("skipme/e.src", "x")   // pruned by skipDir

	isSource := func(p string) bool { return strings.HasSuffix(p, ".src") }
	skipDir := func(p string) error {
		if filepath.Base(p) == "skipme" {
			return filepath.SkipDir
		}
		return nil
	}

	files, err := collectSourceFiles(root, isSource, skipDir)
	if err != nil {
		t.Fatalf("collectSourceFiles: %v", err)
	}
	var rels []string
	for _, f := range files {
		rels = append(rels, f.rel)
		if !filepath.IsAbs(f.abs) {
			t.Errorf("abs not absolute: %q", f.abs)
		}
	}
	sort.Strings(rels)
	want := []string{"a.src", filepath.Join("sub", "d.src")}
	if strings.Join(rels, ",") != strings.Join(want, ",") {
		t.Errorf("collected %v, want %v", rels, want)
	}
}

func TestFallbackFiles(t *testing.T) {
	calls := 0
	idx := &TreeSitterIndexer{
		fileFallback: func(rel, abs string) []model.Function {
			calls++
			return []model.Function{{Name: rel}}
		},
	}
	got := idx.fallbackFiles([]sourceFile{
		{rel: "a.ts", abs: "/x/a.ts"},
		{rel: "b.ts", abs: "/x/b.ts"},
	})
	if calls != 2 || len(got) != 2 {
		t.Fatalf("fallbackFiles calls=%d funcs=%d", calls, len(got))
	}
}
