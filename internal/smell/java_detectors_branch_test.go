package smell

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// jName builds a `name` field identifier leaf.
func jName(text string) *treesitter.Node {
	return &treesitter.Node{Type: "identifier", Field: "name", Text: text}
}

// jInvocation builds a method_invocation node with the given name and (optional)
// argument_list, on the given 0-based start row.
func jInvocation(name string, sRow int, args *treesitter.Node) *treesitter.Node {
	n := &treesitter.Node{Type: "method_invocation", SRow: sRow, ERow: sRow}
	if name != "" {
		n.Children = append(n.Children, jName(name))
	}
	if args != nil {
		n.Children = append(n.Children, args)
	}
	return n
}

// jArgs builds an argument_list with the given literal child types (e.g. "true").
func jArgs(childTypes ...string) *treesitter.Node {
	a := &treesitter.Node{Type: "argument_list", Field: "arguments"}
	for _, ct := range childTypes {
		a.Children = append(a.Children, &treesitter.Node{Type: ct})
	}
	return a
}

func TestJavaArgsHaveTrue(t *testing.T) {
	if javaArgsHaveTrue(nil) {
		t.Error("nil args should be false")
	}
	if !javaArgsHaveTrue(jArgs("true")) {
		t.Error("args with true literal should be true")
	}
	if javaArgsHaveTrue(jArgs("false")) {
		t.Error("args with only false literal should be false")
	}
	if javaArgsHaveTrue(jArgs()) {
		t.Error("empty args should be false")
	}
}

func TestDetectJavaReflect(t *testing.T) {
	root := &treesitter.Node{Type: "program", Children: []*treesitter.Node{
		jInvocation("getDeclaredMethod", 1, jArgs()),
		jInvocation("getDeclaredField", 2, jArgs()),
		jInvocation("getDeclaredMethods", 3, jArgs()),
		jInvocation("getDeclaredFields", 4, jArgs()),
		jInvocation("ordinaryCall", 5, jArgs()), // not reflective -> no finding
		jInvocation("", 6, jArgs()),             // no name field -> no finding
		{Type: "expression_statement", SRow: 7}, // non-invocation node -> skipped
	}}
	findings := detectJavaReflect(root, "X.java")
	if len(findings) != 4 {
		t.Fatalf("got %d findings, want 4: %+v", len(findings), findings)
	}
	for _, f := range findings {
		if f.Rule != "TS-REFL-JV-001" || f.File != "X.java" {
			t.Errorf("unexpected finding: %+v", f)
		}
	}
	// Line is the node's 1-based start line.
	if findings[0].Line != 2 {
		t.Errorf("first finding line = %d, want 2", findings[0].Line)
	}
}

func TestDetectJavaSetAccessible(t *testing.T) {
	root := &treesitter.Node{Type: "program", Children: []*treesitter.Node{
		jInvocation("setAccessible", 1, jArgs("true")),  // fires
		jInvocation("setAccessible", 2, jArgs("false")), // restore guard -> no
		jInvocation("setAccessible", 3, nil),            // no args node -> no
		jInvocation("otherMethod", 4, jArgs("true")),    // wrong name -> no
		jInvocation("", 5, jArgs("true")),               // no name field -> no
		{Type: "expression_statement", SRow: 6},         // non-invocation -> skipped
	}}
	findings := detectJavaSetAccessible(root, "X.java")
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %+v", len(findings), findings)
	}
	f := findings[0]
	if f.Rule != "TS-REFL-JV-002" || f.File != "X.java" || f.Line != 2 || f.Note != "setAccessible(true)" {
		t.Errorf("unexpected finding: %+v", f)
	}
}
