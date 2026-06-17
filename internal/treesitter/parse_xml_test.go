package treesitter

import "testing"

// canned `tree-sitter parse --xml` output for `function foo() {}` — exercises the
// XML→tree marshalling with no CLI dependency (always-green unit coverage).
const cannedFuncXML = `<?xml version="1.0"?>
<sources>
  <source name="/x/foo.ts">
    <program srow="0" scol="0" erow="3" ecol="0">
      <function_declaration srow="0" scol="0" erow="2" ecol="1">
        function
        <identifier field="name" srow="0" scol="9" erow="0" ecol="12">foo</identifier>
        <statement_block srow="0" scol="14" erow="2" ecol="1"></statement_block>
      </function_declaration>
    </program>
  </source>
</sources>`

func TestParseXMLBuildsTree(t *testing.T) {
	sources, err := ParseXML([]byte(cannedFuncXML))
	if err != nil {
		t.Fatalf("ParseXML error: %v", err)
	}
	if len(sources) != 1 {
		t.Fatalf("want 1 source, got %d", len(sources))
	}
	src := sources[0]
	if src.Name != "/x/foo.ts" {
		t.Errorf("source name = %q", src.Name)
	}
	if src.Root == nil || src.Root.Type != "program" {
		t.Fatalf("root = %+v", src.Root)
	}
	fn := src.Root.ChildByType("function_declaration")
	if fn == nil {
		t.Fatal("no function_declaration child")
	}
	if got := fn.StartLine(); got != 1 {
		t.Errorf("StartLine = %d, want 1", got)
	}
	if got := fn.EndLine(); got != 3 {
		t.Errorf("EndLine = %d, want 3", got)
	}
	name := fn.ChildByField("name")
	if name == nil || name.Text != "foo" {
		t.Errorf("name node = %+v", name)
	}
}

func TestParseXMLRejectsGarbage(t *testing.T) {
	if _, err := ParseXML([]byte("<sources><source name=\"x\"><program")); err == nil {
		t.Error("expected error on truncated XML")
	}
}
