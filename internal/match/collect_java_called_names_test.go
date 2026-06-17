package match

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// cannedJavaCallXML mirrors `tree-sitter parse --xml` for a JUnit test body. It
// exercises every branch of collectJavaCalledNames:
//   - method_invocation with a non-empty `name` field (collected)
//   - method_invocation with an empty-text `name` field (skipped)
//   - method_invocation with no `name` field at all (skipped)
//   - object_creation_expression of a bare type_identifier (collected)
//   - object_creation_expression of a generic_type (first type_identifier)
//   - object_creation_expression whose type yields no name (skipped)
const cannedJavaCallXML = `<?xml version="1.0"?>
<sources>
  <source name="/t/FooTest.java">
    <program srow="0" scol="0" erow="9" ecol="0">
      <expression_statement srow="0" scol="0" erow="0" ecol="9">
        <method_invocation srow="0" scol="0" erow="0" ecol="8">
          <identifier field="object" srow="0" scol="0" erow="0" ecol="3">obj</identifier>
          <identifier field="name" srow="0" scol="4" erow="0" ecol="8">doIt</identifier>
          <argument_list field="arguments" srow="0" scol="8" erow="0" ecol="10"></argument_list>
        </method_invocation>
      </expression_statement>
      <expression_statement srow="1" scol="0" erow="1" ecol="9">
        <method_invocation srow="1" scol="0" erow="1" ecol="8">
          <field_access field="object" srow="1" scol="0" erow="1" ecol="3">x.y</field_access>
          <argument_list field="arguments" srow="1" scol="8" erow="1" ecol="10"></argument_list>
        </method_invocation>
      </expression_statement>
      <expression_statement srow="2" scol="0" erow="2" ecol="9">
        <method_invocation srow="2" scol="0" erow="2" ecol="8">
          <identifier field="name" srow="2" scol="4" erow="2" ecol="4"></identifier>
          <argument_list field="arguments" srow="2" scol="8" erow="2" ecol="10"></argument_list>
        </method_invocation>
      </expression_statement>
      <expression_statement srow="3" scol="0" erow="3" ecol="14">
        <object_creation_expression srow="3" scol="0" erow="3" ecol="13">
          new
          <type_identifier field="type" srow="3" scol="4" erow="3" ecol="10">Widget</type_identifier>
          <argument_list field="arguments" srow="3" scol="10" erow="3" ecol="12"></argument_list>
        </object_creation_expression>
      </expression_statement>
      <expression_statement srow="4" scol="0" erow="4" ecol="20">
        <object_creation_expression srow="4" scol="0" erow="4" ecol="19">
          new
          <generic_type field="type" srow="4" scol="4" erow="4" ecol="16">
            <type_identifier srow="4" scol="4" erow="4" ecol="7">Box</type_identifier>
            <type_arguments srow="4" scol="7" erow="4" ecol="16">
              <type_identifier srow="4" scol="8" erow="4" ecol="12">Item</type_identifier>
            </type_arguments>
          </generic_type>
          <argument_list field="arguments" srow="4" scol="16" erow="4" ecol="18"></argument_list>
        </object_creation_expression>
      </expression_statement>
      <expression_statement srow="5" scol="0" erow="5" ecol="14">
        <object_creation_expression srow="5" scol="0" erow="5" ecol="13">
          new
          <void_type field="type" srow="5" scol="4" erow="5" ecol="8">void</void_type>
          <argument_list field="arguments" srow="5" scol="8" erow="5" ecol="10"></argument_list>
        </object_creation_expression>
      </expression_statement>
    </program>
  </source>
</sources>`

func TestCollectJavaCalledNames(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedJavaCallXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	names := collectJavaCalledNames(sources[0].Root)

	for _, want := range []string{"doIt", "Widget", "Box"} {
		if _, ok := names[want]; !ok {
			t.Errorf("missing called name %q in %v", want, names)
		}
	}
	for _, notWant := range []string{"", "Item", "void"} {
		if _, ok := names[notWant]; ok {
			t.Errorf("unexpected name %q collected in %v", notWant, names)
		}
	}
}

func TestJavaConstructedTypeName(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want string
	}{
		{"nil node", nil, ""},
		{"bare type_identifier", &treesitter.Node{Type: "type_identifier", Text: "Foo"}, "Foo"},
		{
			desc: "generic_type uses first type_identifier",
			node: &treesitter.Node{Type: "generic_type", Children: []*treesitter.Node{
				{Type: "type_identifier", Text: "Map"},
				{Type: "type_arguments", Children: []*treesitter.Node{
					{Type: "type_identifier", Text: "K"},
					{Type: "type_identifier", Text: "V"},
				}},
			}},
			want: "Map",
		},
		{
			desc: "empty-text bare type_identifier falls through to walk",
			node: &treesitter.Node{Type: "type_identifier", Text: "", Children: []*treesitter.Node{
				{Type: "type_identifier", Text: "Inner"},
			}},
			want: "Inner",
		},
		{
			desc: "no type_identifier present",
			node: &treesitter.Node{Type: "void_type", Text: "void"},
			want: "",
		},
	}
	for _, c := range cases {
		if got := javaConstructedTypeName(c.node); got != c.want {
			t.Errorf("%s: javaConstructedTypeName = %q, want %q", c.desc, got, c.want)
		}
	}
}

func TestJavaRefsToTestMatch(t *testing.T) {
	tm, ok := javaRefsToTestMatch([]string{"A.java", "A.java", "B.java"})
	if !ok {
		t.Fatal("want found")
	}
	if len(tm.Files) != 2 || tm.Files[0] != "A.java" || tm.Files[1] != "B.java" {
		t.Errorf("Files = %v, want [A.java B.java]", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil", tm.TestFuncs)
	}
	if _, ok := javaRefsToTestMatch(nil); ok {
		t.Error("empty refs should report not found")
	}
}
