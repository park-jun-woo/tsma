package match

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// TestCsSimpleTypeOrName covers every branch of the shared simple-name resolver:
// nil, bare identifier, qualified_name (trailing segment), generic_name (head),
// default-with-text, and default-empty.
func TestCsSimpleTypeOrName(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want string
	}{
		{"nil node", nil, ""},
		{"bare identifier", &treesitter.Node{Type: "identifier", Text: "Foo"}, "Foo"},
		{
			desc: "qualified_name uses name field trailing segment",
			node: &treesitter.Node{Type: "qualified_name", Children: []*treesitter.Node{
				{Type: "identifier", Text: "Ns"},
				{Type: "identifier", Field: "name", Text: "Foo"},
			}},
			want: "Foo",
		},
		{
			desc: "generic_name uses head identifier, never a type argument",
			node: &treesitter.Node{Type: "generic_name", Children: []*treesitter.Node{
				{Type: "identifier", Text: "Box"},
				{Type: "type_argument_list", Children: []*treesitter.Node{
					{Type: "identifier", Text: "Item"},
				}},
			}},
			want: "Box",
		},
		{
			desc: "default node with text",
			node: &treesitter.Node{Type: "predefined_type", Text: "string"},
			want: "string",
		},
		{
			desc: "default node without text",
			node: &treesitter.Node{Type: "predefined_type", Text: ""},
			want: "",
		},
	}
	for _, c := range cases {
		if got := csSimpleTypeOrName(c.node); got != c.want {
			t.Errorf("%s: csSimpleTypeOrName = %q, want %q", c.desc, got, c.want)
		}
	}
}

// TestCsGenericHead covers the identifier-head, qualified_name-head, and
// none-present branches.
func TestCsGenericHead(t *testing.T) {
	identHead := &treesitter.Node{Type: "generic_name", Children: []*treesitter.Node{
		{Type: "type_argument_list"},
		{Type: "identifier", Text: "Box"},
	}}
	if got := csGenericHead(identHead); got == nil || got.Text != "Box" {
		t.Errorf("identifier head: csGenericHead = %+v, want Box", got)
	}

	qualHead := &treesitter.Node{Type: "generic_name", Children: []*treesitter.Node{
		{Type: "qualified_name", Children: []*treesitter.Node{{Type: "identifier", Text: "Ns"}}},
	}}
	if got := csGenericHead(qualHead); got == nil || got.Type != "qualified_name" {
		t.Errorf("qualified head: csGenericHead = %+v, want qualified_name", got)
	}

	none := &treesitter.Node{Type: "generic_name", Children: []*treesitter.Node{
		{Type: "type_argument_list"},
	}}
	if got := csGenericHead(none); got != nil {
		t.Errorf("no head: csGenericHead = %+v, want nil", got)
	}
}

// TestCsInvokedName covers the nil, member_access_expression (name field), and
// default (bare identifier / generic) branches.
func TestCsInvokedName(t *testing.T) {
	if got := csInvokedName(nil); got != "" {
		t.Errorf("nil: csInvokedName = %q, want empty", got)
	}

	// member_access_expression -> name field.
	member := &treesitter.Node{Type: "member_access_expression", Children: []*treesitter.Node{
		{Type: "identifier", Text: "obj"},
		{Type: "identifier", Field: "name", Text: "Foo"},
	}}
	if got := csInvokedName(member); got != "Foo" {
		t.Errorf("member access: csInvokedName = %q, want Foo", got)
	}

	// bare identifier default branch.
	bare := &treesitter.Node{Type: "identifier", Text: "Bar"}
	if got := csInvokedName(bare); got != "Bar" {
		t.Errorf("bare: csInvokedName = %q, want Bar", got)
	}

	// generic_name default branch.
	generic := &treesitter.Node{Type: "generic_name", Children: []*treesitter.Node{
		{Type: "identifier", Text: "Gen"},
		{Type: "type_argument_list"},
	}}
	if got := csInvokedName(generic); got != "Gen" {
		t.Errorf("generic: csInvokedName = %q, want Gen", got)
	}
}

// TestCsConstructedTypeName covers delegation to csSimpleTypeOrName: bare,
// generic, qualified, and nil.
func TestCsConstructedTypeName(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want string
	}{
		{"nil", nil, ""},
		{"bare identifier", &treesitter.Node{Type: "identifier", Text: "Foo"}, "Foo"},
		{
			desc: "generic_name head",
			node: &treesitter.Node{Type: "generic_name", Children: []*treesitter.Node{
				{Type: "identifier", Text: "List"},
				{Type: "type_argument_list"},
			}},
			want: "List",
		},
		{
			desc: "qualified_name trailing",
			node: &treesitter.Node{Type: "qualified_name", Children: []*treesitter.Node{
				{Type: "identifier", Text: "Ns"},
				{Type: "identifier", Field: "name", Text: "Foo"},
			}},
			want: "Foo",
		},
	}
	for _, c := range cases {
		if got := csConstructedTypeName(c.node); got != c.want {
			t.Errorf("%s: csConstructedTypeName = %q, want %q", c.desc, got, c.want)
		}
	}
}

// TestCsRefsToTestMatch covers dedup (order-preserving) and the empty-input
// not-found branch.
func TestCsRefsToTestMatch(t *testing.T) {
	tm, ok := csRefsToTestMatch([]string{"A.cs", "A.cs", "B.cs"})
	if !ok {
		t.Fatal("want found")
	}
	if len(tm.Files) != 2 || tm.Files[0] != "A.cs" || tm.Files[1] != "B.cs" {
		t.Errorf("Files = %v, want [A.cs B.cs]", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil", tm.TestFuncs)
	}
	if _, ok := csRefsToTestMatch(nil); ok {
		t.Error("empty refs should report not found")
	}
}

// TestIsCsTestFile covers the non-.cs, Test-suffix, Tests-suffix, and
// non-test-stem branches.
func TestIsCsTestFile(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"FooTests.cs", true},
		{"FooTest.cs", true},
		{"Foo.cs", false},
		{"FooTests.txt", false},
		{"Tests.cs", true},
		{"Helper.java", false},
	}
	for _, c := range cases {
		if got := isCsTestFile(c.name); got != c.want {
			t.Errorf("isCsTestFile(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}

// cannedCsCallXML mirrors `tree-sitter parse --xml` for an xUnit test body,
// exercising every branch of collectCsCalledNames / csInvokedName /
// csConstructedTypeName:
//   - invocation via member_access_expression name field (obj.Foo())   -> Foo
//   - invocation of a bare identifier (Bar())                          -> Bar
//   - invocation of a generic_name (Gen<T>())                          -> Gen
//   - invocation whose function yields no name (member access, no name) -> skipped
//   - object_creation_expression of a bare identifier (new Made())     -> Made
//   - object_creation_expression of a generic_name (new GenType<T>())  -> GenType
//   - object_creation_expression whose type yields nothing             -> skipped
const cannedCsCallXML = `<?xml version="1.0"?>
<sources>
  <source name="/t/FooTests.cs">
    <compilation_unit srow="0" scol="0" erow="9" ecol="0">
      <invocation_expression srow="0" scol="0" erow="0" ecol="8">
        <member_access_expression field="function" srow="0" scol="0" erow="0" ecol="6">
          <identifier srow="0" scol="0" erow="0" ecol="3">obj</identifier>
          <identifier field="name" srow="0" scol="4" erow="0" ecol="6">Foo</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="0" scol="6" erow="0" ecol="8"></argument_list>
      </invocation_expression>
      <invocation_expression srow="1" scol="0" erow="1" ecol="5">
        <identifier field="function" srow="1" scol="0" erow="1" ecol="3">Bar</identifier>
        <argument_list field="arguments" srow="1" scol="3" erow="1" ecol="5"></argument_list>
      </invocation_expression>
      <invocation_expression srow="2" scol="0" erow="2" ecol="8">
        <generic_name field="function" srow="2" scol="0" erow="2" ecol="6">
          <identifier srow="2" scol="0" erow="2" ecol="3">Gen</identifier>
          <type_argument_list srow="2" scol="3" erow="2" ecol="6">
            <identifier srow="2" scol="4" erow="2" ecol="5">T</identifier>
          </type_argument_list>
        </generic_name>
        <argument_list field="arguments" srow="2" scol="6" erow="2" ecol="8"></argument_list>
      </invocation_expression>
      <invocation_expression srow="3" scol="0" erow="3" ecol="8">
        <member_access_expression field="function" srow="3" scol="0" erow="3" ecol="6">
          <identifier srow="3" scol="0" erow="3" ecol="3">obj</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="3" scol="6" erow="3" ecol="8"></argument_list>
      </invocation_expression>
      <object_creation_expression srow="4" scol="0" erow="4" ecol="10">
        <identifier field="type" srow="4" scol="4" erow="4" ecol="8">Made</identifier>
        <argument_list field="arguments" srow="4" scol="8" erow="4" ecol="10"></argument_list>
      </object_creation_expression>
      <object_creation_expression srow="5" scol="0" erow="5" ecol="14">
        <generic_name field="type" srow="5" scol="4" erow="5" ecol="12">
          <identifier srow="5" scol="4" erow="5" ecol="11">GenType</identifier>
          <type_argument_list srow="5" scol="11" erow="5" ecol="12">
            <identifier srow="5" scol="12" erow="5" ecol="13">T</identifier>
          </type_argument_list>
        </generic_name>
        <argument_list field="arguments" srow="5" scol="12" erow="5" ecol="14"></argument_list>
      </object_creation_expression>
      <object_creation_expression srow="6" scol="0" erow="6" ecol="10">
        <predefined_type field="type" srow="6" scol="4" erow="6" ecol="8"></predefined_type>
        <argument_list field="arguments" srow="6" scol="8" erow="6" ecol="10"></argument_list>
      </object_creation_expression>
    </compilation_unit>
  </source>
</sources>`

// TestCollectCsCalledNames exercises the Walk + switch over invocation and
// object-creation nodes using canned XML (no tree-sitter CLI needed).
func TestCollectCsCalledNames(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedCsCallXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	names := collectCsCalledNames(sources[0].Root)

	for _, want := range []string{"Foo", "Bar", "Gen", "Made", "GenType"} {
		if _, ok := names[want]; !ok {
			t.Errorf("missing called name %q in %v", want, names)
		}
	}
	for _, notWant := range []string{"", "obj", "T"} {
		if _, ok := names[notWant]; ok {
			t.Errorf("unexpected name %q collected in %v", notWant, names)
		}
	}
}

// TestCsTestDirs covers the same-dir-only ("", "."), single-segment, and
// multi-segment (parallel *.Tests project) branches.
func TestCsTestDirs(t *testing.T) {
	cases := []struct {
		src  string
		want []string
	}{
		{"", []string{""}},
		{".", []string{"."}},
		{"App", []string{"App", "App.Tests"}},
		{"App/Services", []string{"App/Services", "App.Tests/Services"}},
		{"App\\Services", []string{"App\\Services", "App.Tests/Services"}},
	}
	for _, c := range cases {
		got := csTestDirs(c.src)
		if len(got) != len(c.want) {
			t.Errorf("csTestDirs(%q) = %v, want %v", c.src, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("csTestDirs(%q)[%d] = %q, want %q", c.src, i, got[i], c.want[i])
			}
		}
	}
}
