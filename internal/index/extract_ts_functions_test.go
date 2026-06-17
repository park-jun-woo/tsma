package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// cannedTSXML mirrors real `tree-sitter parse --xml` shapes: an exported
// function, an exported const arrow, a non-exported function, and a class with
// an exported method + a constructor (which must be skipped). It lets the TS
// extractor be unit-tested with no tree-sitter CLI present.
const cannedTSXML = `<?xml version="1.0"?>
<sources>
  <source name="/p/svc/api.ts">
    <program srow="0" scol="0" erow="20" ecol="0">
      <export_statement srow="0" scol="0" erow="3" ecol="1">
        export
        <function_declaration field="declaration" srow="0" scol="7" erow="3" ecol="1">
          function
          <identifier field="name" srow="0" scol="16" erow="0" ecol="19">add</identifier>
        </function_declaration>
      </export_statement>
      <export_statement srow="5" scol="0" erow="5" ecol="40">
        export
        <lexical_declaration srow="5" scol="7" erow="5" ecol="40">
          const
          <variable_declarator srow="5" scol="13" erow="5" ecol="39">
            <identifier field="name" srow="5" scol="13" erow="5" ecol="21">classify</identifier>
            <arrow_function field="value" srow="5" scol="24" erow="5" ecol="39"></arrow_function>
          </variable_declarator>
        </lexical_declaration>
      </export_statement>
      <function_declaration srow="7" scol="0" erow="9" ecol="1">
        <identifier field="name" srow="7" scol="9" erow="7" ecol="23">internalHelper</identifier>
      </function_declaration>
      <export_statement srow="11" scol="0" erow="18" ecol="1">
        <class_declaration field="declaration" srow="11" scol="7" erow="18" ecol="1">
          <identifier field="name" srow="11" scol="13" erow="11" ecol="22">Rectangle</identifier>
          <class_body srow="11" scol="23" erow="18" ecol="1">
            <method_definition srow="12" scol="2" erow="12" ecol="30">
              <property_identifier field="name" srow="12" scol="2" erow="12" ecol="13">constructor</property_identifier>
            </method_definition>
            <method_definition srow="14" scol="2" erow="16" ecol="3">
              <property_identifier field="name" srow="14" scol="2" erow="14" ecol="6">Area</property_identifier>
            </method_definition>
          </class_body>
        </class_declaration>
      </export_statement>
    </program>
  </source>
</sources>`

func TestExtractTSFunctions(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedTSXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	relPath := "svc/api.ts"
	funcs := extractTSFunctions(sources[0].Root, relPath, pkgDirOf(relPath))

	type want struct {
		qn       string
		start    int
		end      int
		exported bool
	}
	wants := map[string]want{
		"add":            {"svc.add", 1, 4, true},
		"classify":       {"svc.classify", 6, 6, true},
		"internalHelper": {"svc.internalHelper", 8, 10, false},
		"Area":           {"svc.Rectangle.Area", 15, 17, true},
	}

	if len(funcs) != len(wants) {
		t.Fatalf("got %d funcs, want %d: %+v", len(funcs), len(wants), funcs)
	}
	for _, f := range funcs {
		w, ok := wants[f.Name]
		if !ok {
			t.Errorf("unexpected func %q (constructor must be skipped)", f.Name)
			continue
		}
		if f.QualifiedName != w.qn {
			t.Errorf("%s: QualifiedName = %q, want %q", f.Name, f.QualifiedName, w.qn)
		}
		if f.StartLine != w.start || f.EndLine != w.end {
			t.Errorf("%s: range = %d..%d, want %d..%d", f.Name, f.StartLine, f.EndLine, w.start, w.end)
		}
		if f.Exported != w.exported {
			t.Errorf("%s: Exported = %v, want %v", f.Name, f.Exported, w.exported)
		}
		if f.File != relPath {
			t.Errorf("%s: File = %q", f.Name, f.File)
		}
	}
}
