package index

import "testing"

func TestProcessRsLineFn(t *testing.T) {
	st := &rsParseState{relPath: "lib.rs"}
	processRsLine(st, "pub fn add() {")
	processRsLine(st, "    a + b")
	processRsLine(st, "}")

	if len(st.functions) != 1 || st.functions[0].Name != "add" {
		t.Fatalf("functions = %+v, want one add", st.functions)
	}
	if st.depth != 0 {
		t.Errorf("depth = %d, want 0 after closing brace", st.depth)
	}
}

func TestProcessRsLineCfgTestCarriesOver(t *testing.T) {
	st := &rsParseState{relPath: "lib.rs"}
	processRsLine(st, "#[cfg(test)]")
	processRsLine(st, "mod tests {")
	processRsLine(st, "fn it_works() {}")

	if len(st.functions) != 0 {
		t.Errorf("cfg(test) module fns should be skipped, got %+v", st.functions)
	}
}
