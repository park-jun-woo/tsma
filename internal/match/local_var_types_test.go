package match

import "testing"

func TestLocalVarTypesSingleBinding(t *testing.T) {
	body := parseFuncBody(t, `package p
func f() {
	a := &GoFile{}
	b := CSharpFile{}
	_ = a
	_ = b
}
`, "f")
	got := localVarTypes(body)
	if got["a"] != "GoFile" {
		t.Errorf("a = %q, want GoFile", got["a"])
	}
	if got["b"] != "CSharpFile" {
		t.Errorf("b = %q, want CSharpFile", got["b"])
	}
}

func TestLocalVarTypesReassignedDropped(t *testing.T) {
	// A reassigned variable is ambiguous -> dropped (conservative).
	body := parseFuncBody(t, `package p
func f() {
	x := &GoFile{}
	x = &CSharpFile{}
	_ = x
}
`, "f")
	got := localVarTypes(body)
	if _, ok := got["x"]; ok {
		t.Errorf("reassigned x should be dropped, got %q", got["x"])
	}
}

func TestLocalVarTypesNonCompositeDropped(t *testing.T) {
	// RHS is not a composite literal -> not tracked.
	body := parseFuncBody(t, `package p
func f() {
	c := NewCSharpFile()
	d := 5
	_ = c
	_ = d
}
`, "f")
	got := localVarTypes(body)
	if _, ok := got["c"]; ok {
		t.Errorf("constructor-bound c must not be tracked, got %q", got["c"])
	}
	if _, ok := got["d"]; ok {
		t.Errorf("int-bound d must not be tracked, got %q", got["d"])
	}
}

func TestLocalVarTypesMultiAssignDropped(t *testing.T) {
	// Tuple/multi assign poisons the names.
	body := parseFuncBody(t, `package p
func f() {
	a, b := &GoFile{}, &CSharpFile{}
	_ = a
	_ = b
}
`, "f")
	got := localVarTypes(body)
	if _, ok := got["a"]; ok {
		t.Errorf("multi-assign a must not be tracked, got %q", got["a"])
	}
	if _, ok := got["b"]; ok {
		t.Errorf("multi-assign b must not be tracked, got %q", got["b"])
	}
}

func TestLocalVarTypesNonIdentLHSIgnored(t *testing.T) {
	// A single-LHS assignment whose LHS is not a bare Ident (a field/selector
	// assignment) is not a local-variable binding and is skipped without
	// poisoning the local var that backs it. The earlier `g := &GoFile{}` stays
	// tracked; the field write `g.field = &CSharpFile{}` must not retype `g`.
	body := parseFuncBody(t, `package p
func f() {
	g := &GoFile{}
	g.field = &CSharpFile{}
	_ = g
}
`, "f")
	got := localVarTypes(body)
	if got["g"] != "GoFile" {
		t.Errorf("g = %q, want GoFile (field write must not retype the var)", got["g"])
	}
}

func TestLocalVarTypesPoisonedThenRebindStaysDropped(t *testing.T) {
	// A name first poisoned by a multi-assign is removed from `types` (so the
	// "already seen" check does not fire) but recorded as poisoned; a later
	// single `:=` rebind to a composite literal must NOT restore a type. This is
	// the poisoned-guard branch protecting against mis-attribution after the
	// name became ambiguous.
	body := parseFuncBody(t, `package p
func f() {
	x, y := load()      // multi-assign poisons x and y
	x := &GoFile{}      // rebind: poisoned -> must stay unknown
	_ = x
	_ = y
}
`, "f")
	got := localVarTypes(body)
	if _, ok := got["x"]; ok {
		t.Errorf("poisoned-then-rebound x must stay dropped, got %q", got["x"])
	}
}

func TestLocalVarTypesReassignmentPoisonsReceiverResolution(t *testing.T) {
	// End-to-end safety check: when a variable is bound to two different types
	// (re-binding via :=), its type becomes ambiguous, so a later method call on
	// it must resolve to an UNKNOWN receiver rather than a stale/wrong type.
	// This is the "prefer non-detection over mis-attribution" guarantee.
	body := parseFuncBody(t, `package p
func f() {
	f := &GoFile{}
	f := &CSharpFile{}
	f.GetFuncs()
}
`, "f")

	// localVarTypes must drop the ambiguous var entirely.
	if _, ok := localVarTypes(body)["f"]; ok {
		t.Fatalf("ambiguous f must be dropped from local var types")
	}

	// And the downstream ref must carry an unknown ("") receiver, not GoFile or
	// CSharpFile, so the matcher cannot mis-attribute the call.
	refs := collectCalledRefs(body)
	r, ok := receiverOf(refs, "GetFuncs")
	if !ok {
		t.Fatal("GetFuncs ref should be present")
	}
	if r != "" {
		t.Errorf("GetFuncs receiver = %q, want unknown (\"\") after ambiguous rebinding", r)
	}
}

func TestLocalVarTypesNilBody(t *testing.T) {
	got := localVarTypes(nil)
	if got == nil || len(got) != 0 {
		t.Fatalf("localVarTypes(nil) = %v, want empty non-nil", got)
	}
}
