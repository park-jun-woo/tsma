package match

import "testing"

// refNames flattens a calledRef set to the set of bare names (dropping
// receiver), so name-only assertions can reuse hasAll.
func refNames(refs map[calledRef]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for r := range refs {
		out[r.Name] = struct{}{}
	}
	return out
}

// receiverOf returns the receiver recorded for the first calledRef with the
// given name, and whether any such ref exists.
func receiverOf(refs map[calledRef]struct{}, name string) (string, bool) {
	for r := range refs {
		if r.Name == name {
			return r.Receiver, true
		}
	}
	return "", false
}

func TestCollectCalledRefsCompositeLitDirect(t *testing.T) {
	body := parseFuncBody(t, `package p
func f() {
	(&CSharpFile{}).GetFuncs()
	GoFile{}.GetFuncs()
	Plain()
}
`, "f")
	refs := collectCalledRefs(body)

	// Two distinct receivers for the same name -> two entries.
	want := map[calledRef]struct{}{
		{Name: "GetFuncs", Receiver: "CSharpFile"}: {},
		{Name: "GetFuncs", Receiver: "GoFile"}:     {},
		{Name: "Plain", Receiver: ""}:              {},
	}
	if len(refs) != len(want) {
		t.Fatalf("refs = %v, want %v", refs, want)
	}
	for w := range want {
		if _, ok := refs[w]; !ok {
			t.Errorf("missing ref %v in %v", w, refs)
		}
	}
}

func TestCollectCalledRefsLocalVar(t *testing.T) {
	body := parseFuncBody(t, `package p
func f() {
	f := &CSharpFile{}
	f.GetFuncs()
	v := JavaFile{}
	v.GetPath()
}
`, "f")
	refs := collectCalledRefs(body)

	if r, ok := receiverOf(refs, "GetFuncs"); !ok || r != "CSharpFile" {
		t.Errorf("GetFuncs receiver = %q ok=%v, want CSharpFile", r, ok)
	}
	if r, ok := receiverOf(refs, "GetPath"); !ok || r != "JavaFile" {
		t.Errorf("GetPath receiver = %q ok=%v, want JavaFile", r, ok)
	}
}

func TestCollectCalledRefsUnknownReceivers(t *testing.T) {
	// Constructor return, interface argument, field access, reassigned var:
	// all unresolvable -> unknown ("") receiver. No constructor inference.
	body := parseFuncBody(t, `package p
func f(iface SourceFile) {
	NewCSharpFile().GetFuncs()
	iface.GetPath()
	obj.field.GetLang()
	x := &GoFile{}
	x = &CSharpFile{}
	x.GetTypes()
}
`, "f")
	refs := collectCalledRefs(body)

	for _, name := range []string{"GetFuncs", "GetPath", "GetLang", "GetTypes"} {
		r, ok := receiverOf(refs, name)
		if !ok {
			t.Errorf("%s should be present", name)
			continue
		}
		if r != "" {
			t.Errorf("%s receiver = %q, want unknown (\"\")", name, r)
		}
	}
}

func TestCollectCalledRefsEmptyCalleeNameSkipped(t *testing.T) {
	// A callee with no bare identifier name (chained call result, or an
	// immediately-invoked function literal) yields calleeName == "" and must be
	// skipped, never producing a calledRef with an empty Name. The named calls
	// nested inside (Outer, Plain) are still collected so the walk continues.
	body := parseFuncBody(t, `package p
func f() {
	Outer()()            // chained call: outer .Fun is a CallExpr -> name ""
	func() { Plain() }() // IIFE: .Fun is a FuncLit -> name ""
}
`, "f")
	refs := collectCalledRefs(body)

	// No ref with an empty name must ever be recorded.
	if _, bad := receiverOf(refs, ""); bad {
		t.Errorf("empty-name callee must be skipped, got ref %v", refs)
	}
	// The inner named calls survive the walk.
	names := refNames(refs)
	for _, want := range []string{"Outer", "Plain"} {
		if _, ok := names[want]; !ok {
			t.Errorf("expected %s to be collected, got %v", want, names)
		}
	}
}

func TestCollectCalledRefsNilBody(t *testing.T) {
	refs := collectCalledRefs(nil)
	if refs == nil || len(refs) != 0 {
		t.Fatalf("collectCalledRefs(nil) = %v, want empty non-nil", refs)
	}
}
