package smell

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// cannedCsReflectXML mirrors `tree-sitter parse --xml` for a C# test body,
// exercising every branch of detectCsReflect:
//   - invocation via member_access name GetMethod          -> fires
//   - invocation via member_access name GetProperties      -> fires (plural)
//   - invocation via member_access name Foo (not reflect)  -> no
//   - invocation via member_access with no name field      -> no (name==nil)
//   - invocation of a bare identifier (not member_access)  -> no
const cannedCsReflectXML = `<?xml version="1.0"?>
<sources>
  <source name="/t/T.cs">
    <compilation_unit srow="0" scol="0" erow="9" ecol="0">
      <invocation_expression srow="2" scol="0" erow="2" ecol="20">
        <member_access_expression field="function" srow="2" scol="0" erow="2" ecol="12">
          <identifier srow="2" scol="0" erow="2" ecol="1">t</identifier>
          <identifier field="name" srow="2" scol="2" erow="2" ecol="11">GetMethod</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="2" scol="12" erow="2" ecol="20"></argument_list>
      </invocation_expression>
      <invocation_expression srow="3" scol="0" erow="3" ecol="20">
        <member_access_expression field="function" srow="3" scol="0" erow="3" ecol="16">
          <identifier srow="3" scol="0" erow="3" ecol="1">t</identifier>
          <identifier field="name" srow="3" scol="2" erow="3" ecol="15">GetProperties</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="3" scol="16" erow="3" ecol="20"></argument_list>
      </invocation_expression>
      <invocation_expression srow="4" scol="0" erow="4" ecol="10">
        <member_access_expression field="function" srow="4" scol="0" erow="4" ecol="6">
          <identifier srow="4" scol="0" erow="4" ecol="1">t</identifier>
          <identifier field="name" srow="4" scol="2" erow="4" ecol="5">Foo</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="4" scol="6" erow="4" ecol="10"></argument_list>
      </invocation_expression>
      <invocation_expression srow="5" scol="0" erow="5" ecol="10">
        <member_access_expression field="function" srow="5" scol="0" erow="5" ecol="6">
          <identifier srow="5" scol="0" erow="5" ecol="1">t</identifier>
        </member_access_expression>
        <argument_list field="arguments" srow="5" scol="6" erow="5" ecol="10"></argument_list>
      </invocation_expression>
      <invocation_expression srow="6" scol="0" erow="6" ecol="10">
        <identifier field="function" srow="6" scol="0" erow="6" ecol="3">Bar</identifier>
        <argument_list field="arguments" srow="6" scol="3" erow="6" ecol="10"></argument_list>
      </invocation_expression>
    </compilation_unit>
  </source>
</sources>`

// TestDetectCsReflectBranches drives all branches via canned XML and asserts only
// the two reflective calls fire (with the right rule, note, and line).
func TestDetectCsReflectBranches(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedCsReflectXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	findings := detectCsReflect(sources[0].Root, "T.cs")
	if len(findings) != 2 {
		t.Fatalf("findings = %+v, want 2", findings)
	}
	if findings[0].Rule != "TS-REFL-CS-001" || findings[0].Note != "GetMethod()" || findings[0].Line != 3 {
		t.Errorf("findings[0] = %+v", findings[0])
	}
	if findings[1].Note != "GetProperties()" || findings[1].File != "T.cs" {
		t.Errorf("findings[1] = %+v", findings[1])
	}
}

// cannedCsReflectInfoXML exercises every branch of detectCsReflectInfo:
//   - variable_declaration type field MethodInfo    -> fires
//   - variable_declaration type field PropertyInfo  -> fires
//   - variable_declaration type field int (not set) -> no
//   - variable_declaration with no type field       -> no (typeNode==nil)
//   - a non-variable_declaration node               -> no
const cannedCsReflectInfoXML = `<?xml version="1.0"?>
<sources>
  <source name="/t/T.cs">
    <compilation_unit srow="0" scol="0" erow="9" ecol="0">
      <variable_declaration srow="1" scol="0" erow="1" ecol="20">
        <identifier field="type" srow="1" scol="0" erow="1" ecol="10">MethodInfo</identifier>
        <variable_declarator srow="1" scol="11" erow="1" ecol="12">m</variable_declarator>
      </variable_declaration>
      <variable_declaration srow="2" scol="0" erow="2" ecol="20">
        <identifier field="type" srow="2" scol="0" erow="2" ecol="12">PropertyInfo</identifier>
        <variable_declarator srow="2" scol="13" erow="2" ecol="14">p</variable_declarator>
      </variable_declaration>
      <variable_declaration srow="3" scol="0" erow="3" ecol="10">
        <predefined_type field="type" srow="3" scol="0" erow="3" ecol="3">int</predefined_type>
        <variable_declarator srow="3" scol="4" erow="3" ecol="5">x</variable_declarator>
      </variable_declaration>
      <variable_declaration srow="4" scol="0" erow="4" ecol="10">
        <variable_declarator srow="4" scol="4" erow="4" ecol="5">y</variable_declarator>
      </variable_declaration>
      <expression_statement srow="5" scol="0" erow="5" ecol="10">noop</expression_statement>
    </compilation_unit>
  </source>
</sources>`

// TestDetectCsReflectInfoBranches drives all branches via canned XML and asserts
// only the two reflection-handle declarations fire.
func TestDetectCsReflectInfoBranches(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedCsReflectInfoXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	findings := detectCsReflectInfo(sources[0].Root, "T.cs")
	if len(findings) != 2 {
		t.Fatalf("findings = %+v, want 2", findings)
	}
	if findings[0].Rule != "TS-REFL-CS-002" || findings[0].Note != "MethodInfo" || findings[0].Line != 2 {
		t.Errorf("findings[0] = %+v", findings[0])
	}
	if findings[1].Note != "PropertyInfo" {
		t.Errorf("findings[1] = %+v", findings[1])
	}
}

// TestScanCsNoTreeSitter covers the tree-sitter-unavailable branch: ScanCs
// returns (nil, err) so the caller ignores the file.
func TestScanCsNoTreeSitter(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	findings, err := ScanCs("../../testdata/csharp/Calc.Tests/ReflectionTests.cs")
	if err == nil {
		t.Error("expected error when tree-sitter unavailable")
	}
	if findings != nil {
		t.Errorf("expected nil findings, got %+v", findings)
	}
}

// TestScanCsParseError covers the ParseFile-failure branch with a real CLI
// present but a nonexistent source path.
func TestScanCsParseError(t *testing.T) {
	if !locateCsSmell(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	findings, err := ScanCs("../../testdata/csharp/Calc.Tests/does_not_exist.cs")
	if err == nil {
		t.Error("expected parse error for a missing file")
	}
	if findings != nil {
		t.Errorf("expected nil findings, got %+v", findings)
	}
}
