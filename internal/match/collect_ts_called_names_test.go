package match

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// cannedCallXML mirrors `tree-sitter parse --xml` for a test body that calls
// add() and constructs new Rectangle() — lets the call-site collector be tested
// with no tree-sitter CLI present.
const cannedCallXML = `<?xml version="1.0"?>
<sources>
  <source name="/t/x.test.ts">
    <program srow="0" scol="0" erow="3" ecol="0">
      <expression_statement srow="0" scol="0" erow="0" ecol="9">
        <call_expression srow="0" scol="0" erow="0" ecol="8">
          <identifier field="function" srow="0" scol="0" erow="0" ecol="3">add</identifier>
          <arguments field="arguments" srow="0" scol="3" erow="0" ecol="8"></arguments>
        </call_expression>
      </expression_statement>
      <expression_statement srow="1" scol="0" erow="1" ecol="20">
        <new_expression srow="1" scol="0" erow="1" ecol="19">
          new
          <identifier field="constructor" srow="1" scol="4" erow="1" ecol="13">Rectangle</identifier>
          <arguments field="arguments" srow="1" scol="13" erow="1" ecol="19"></arguments>
        </new_expression>
      </expression_statement>
      <expression_statement srow="2" scol="0" erow="2" ecol="15">
        <call_expression srow="2" scol="0" erow="2" ecol="14">
          <member_expression field="function" srow="2" scol="0" erow="2" ecol="8">
            <identifier field="object" srow="2" scol="0" erow="2" ecol="3">obj</identifier>
            <property_identifier field="property" srow="2" scol="4" erow="2" ecol="8">area</property_identifier>
          </member_expression>
          <arguments field="arguments" srow="2" scol="8" erow="2" ecol="14"></arguments>
        </call_expression>
      </expression_statement>
    </program>
  </source>
</sources>`

func TestCollectTSCalledNames(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedCallXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	names := collectTSCalledNames(sources[0].Root)
	for _, want := range []string{"add", "Rectangle", "area"} {
		if _, ok := names[want]; !ok {
			t.Errorf("missing called name %q in %v", want, names)
		}
	}
}

func TestTSRefsToTestMatchDedups(t *testing.T) {
	tm, ok := tsRefsToTestMatch([]string{"a.test.ts", "a.test.ts", "b.test.ts"})
	if !ok {
		t.Fatal("want found")
	}
	if len(tm.Files) != 2 || tm.Files[0] != "a.test.ts" || tm.Files[1] != "b.test.ts" {
		t.Errorf("Files = %v, want [a.test.ts b.test.ts]", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil (run whole file)", tm.TestFuncs)
	}
	if _, ok := tsRefsToTestMatch(nil); ok {
		t.Error("empty refs should report not found")
	}
}
