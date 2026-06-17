package index

import "testing"

// TestParsePyAst_CannedJSON covers the JSON→model.Function mapping without a
// Python interpreter, so the marshaling contract is regression-tested even in
// clean environments. It checks qualified-name composition (pkgDir + receiver),
// the exported heuristic, and that line ranges pass through verbatim.
func TestParsePyAst_CannedJSON(t *testing.T) {
	data := []byte(`[
	  {"name":"classify","start_line":1,"end_line":4,"col":0,"receiver":"","decorators":[]},
	  {"name":"Add","start_line":12,"end_line":15,"col":4,"receiver":"Calculator","decorators":["staticmethod"]},
	  {"name":"_inner","start_line":13,"end_line":14,"col":8,"receiver":"","decorators":[]}
	]`)
	funcs, err := parsePyAst(data, "src/calc.py")
	if err != nil {
		t.Fatalf("parsePyAst: %v", err)
	}
	if len(funcs) != 3 {
		t.Fatalf("got %d functions, want 3", len(funcs))
	}

	c := funcs[0]
	if c.QualifiedName != "src.classify" {
		t.Errorf("classify qualified = %q, want src.classify", c.QualifiedName)
	}
	if c.Exported {
		t.Errorf("classify should be unexported (lowercase)")
	}
	if c.StartLine != 1 || c.EndLine != 4 {
		t.Errorf("classify range = %d-%d, want 1-4", c.StartLine, c.EndLine)
	}

	a := funcs[1]
	if a.QualifiedName != "src.Calculator.Add" {
		t.Errorf("Add qualified = %q, want src.Calculator.Add", a.QualifiedName)
	}
	if a.Receiver != "Calculator" {
		t.Errorf("Add receiver = %q, want Calculator", a.Receiver)
	}
	if !a.Exported {
		t.Errorf("Add should be exported (uppercase, no leading underscore)")
	}

	in := funcs[2]
	if in.Receiver != "" {
		t.Errorf("_inner receiver = %q, want empty", in.Receiver)
	}
	if in.Exported {
		t.Errorf("_inner should be unexported (leading underscore)")
	}
}

// TestParsePyAst_Invalid ensures malformed JSON surfaces an error (the caller
// then per-file-falls-back to the line indexer).
func TestParsePyAst_Invalid(t *testing.T) {
	if _, err := parsePyAst([]byte("not json"), "x.py"); err == nil {
		t.Fatal("expected error for malformed ast JSON")
	}
}
