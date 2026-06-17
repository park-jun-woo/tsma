package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// cannedJavaXML mirrors real `tree-sitter parse --xml` shapes for a Java file.
// It exercises every branch of the Java extractor without a tree-sitter CLI:
//   - a dotted package (scoped_identifier) -> javaPackageName joins leaves
//   - public/private/no modifiers -> javaNodeExported true/false/no-modifiers
//   - method, constructor, compact_constructor declarations -> emitted
//   - nameless and empty-named members -> javaFuncFromNode returns false
//   - nested class -> walkJavaTypeDecl pushes a scope
//   - enum + enum_body_declarations wrapper -> descended without a new scope
//   - interface / record / annotation type declarations -> all dispatched
//   - nameless / empty-named / body-less types -> walkJavaTypeDecl early returns
//   - a field_declaration -> dispatchJavaMember default (ignored)
const cannedJavaXML = `<?xml version="1.0"?>
<sources>
  <source name="/p/calc/Calculator.java">
    <program srow="0" scol="0" erow="60" ecol="0">
      <package_declaration srow="0" scol="0" erow="0" ecol="20">
        package
        <scoped_identifier srow="0" scol="8" erow="0" ecol="19">
          <scoped_identifier srow="0" scol="8" erow="0" ecol="15">
            <identifier srow="0" scol="8" erow="0" ecol="11">com</identifier>
            <identifier srow="0" scol="12" erow="0" ecol="15">example</identifier>
          </scoped_identifier>
          <identifier srow="0" scol="16" erow="0" ecol="19">app</identifier>
        </scoped_identifier>
      </package_declaration>
      <class_declaration srow="2" scol="0" erow="40" ecol="1">
        <modifiers srow="2" scol="0" erow="2" ecol="6">public</modifiers>
        <identifier field="name" srow="2" scol="13" erow="2" ecol="23">Calculator</identifier>
        <class_body field="body" srow="2" scol="24" erow="40" ecol="1">
          <method_declaration srow="3" scol="2" erow="5" ecol="3">
            <modifiers srow="3" scol="2" erow="3" ecol="8">public</modifiers>
            <identifier field="name" srow="3" scol="13" erow="3" ecol="16">add</identifier>
          </method_declaration>
          <constructor_declaration srow="6" scol="2" erow="8" ecol="3">
            <modifiers srow="6" scol="2" erow="6" ecol="8">public</modifiers>
            <identifier field="name" srow="6" scol="9" erow="6" ecol="19">Calculator</identifier>
          </constructor_declaration>
          <method_declaration srow="9" scol="2" erow="11" ecol="3">
            <modifiers srow="9" scol="2" erow="9" ecol="16">private static</modifiers>
            <identifier field="name" srow="9" scol="17" erow="9" ecol="23">helper</identifier>
          </method_declaration>
          <method_declaration srow="12" scol="2" erow="14" ecol="3">
            <identifier field="name" srow="12" scol="2" erow="12" ecol="8">noMods</identifier>
          </method_declaration>
          <method_declaration srow="15" scol="2" erow="15" ecol="20"></method_declaration>
          <method_declaration srow="16" scol="2" erow="16" ecol="20">
            <identifier field="name" srow="16" scol="2" erow="16" ecol="2"></identifier>
          </method_declaration>
          <class_declaration srow="17" scol="2" erow="20" ecol="3">
            <identifier field="name" srow="17" scol="8" erow="17" ecol="13">Inner</identifier>
            <class_body field="body" srow="17" scol="14" erow="20" ecol="3">
              <method_declaration srow="18" scol="4" erow="19" ecol="5">
                <identifier field="name" srow="18" scol="4" erow="18" ecol="11">innerM</identifier>
              </method_declaration>
            </class_body>
          </class_declaration>
          <enum_declaration srow="21" scol="2" erow="26" ecol="3">
            <identifier field="name" srow="21" scol="7" erow="21" ecol="12">Color</identifier>
            <enum_body field="body" srow="21" scol="13" erow="26" ecol="3">
              <enum_body_declarations srow="22" scol="4" erow="25" ecol="5">
                <method_declaration srow="23" scol="4" erow="24" ecol="5">
                  <identifier field="name" srow="23" scol="4" erow="23" ecol="14">enumMethod</identifier>
                </method_declaration>
              </enum_body_declarations>
            </enum_body>
          </enum_declaration>
          <field_declaration srow="27" scol="2" erow="27" ecol="20"></field_declaration>
          <class_declaration srow="28" scol="2" erow="28" ecol="10"></class_declaration>
          <class_declaration srow="29" scol="2" erow="29" ecol="20">
            <identifier field="name" srow="29" scol="8" erow="29" ecol="8"></identifier>
          </class_declaration>
          <class_declaration srow="30" scol="2" erow="30" ecol="20">
            <identifier field="name" srow="30" scol="8" erow="30" ecol="14">NoBody</identifier>
          </class_declaration>
        </class_body>
      </class_declaration>
      <interface_declaration srow="41" scol="0" erow="44" ecol="1">
        <identifier field="name" srow="41" scol="10" erow="41" ecol="15">Shape</identifier>
        <interface_body field="body" srow="41" scol="16" erow="44" ecol="1">
          <method_declaration srow="42" scol="2" erow="42" ecol="20">
            <modifiers srow="42" scol="2" erow="42" ecol="8">public abstract</modifiers>
            <identifier field="name" srow="42" scol="9" erow="42" ecol="13">area</identifier>
          </method_declaration>
        </interface_body>
      </interface_declaration>
      <record_declaration srow="45" scol="0" erow="47" ecol="1">
        <identifier field="name" srow="45" scol="7" erow="45" ecol="12">Point</identifier>
        <class_body field="body" srow="45" scol="13" erow="47" ecol="1">
          <compact_constructor_declaration srow="46" scol="2" erow="46" ecol="20">
            <identifier field="name" srow="46" scol="2" erow="46" ecol="7">Point</identifier>
          </compact_constructor_declaration>
        </class_body>
      </record_declaration>
      <annotation_type_declaration srow="48" scol="0" erow="49" ecol="1">
        <identifier field="name" srow="48" scol="0" erow="48" ecol="10">MyAnno</identifier>
        <annotation_type_body field="body" srow="48" scol="11" erow="49" ecol="1"></annotation_type_body>
      </annotation_type_declaration>
    </program>
  </source>
</sources>`

func TestExtractJavaFunctionsCanned(t *testing.T) {
	sources, err := treesitter.ParseXML([]byte(cannedJavaXML))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	relPath := "calc/Calculator.java"
	funcs := extractJavaFunctions(sources[0].Root, relPath, pkgDirOf(relPath))

	type want struct {
		qn         string
		start, end int
		exported   bool
	}
	wants := map[string]want{
		"add":        {"com.example.app.Calculator.add", 4, 6, true},
		"Calculator": {"com.example.app.Calculator.Calculator", 7, 9, true},
		"helper":     {"com.example.app.Calculator.helper", 10, 12, false},
		"noMods":     {"com.example.app.Calculator.noMods", 13, 15, false},
		"innerM":     {"com.example.app.Calculator.Inner.innerM", 19, 20, false},
		"enumMethod": {"com.example.app.Calculator.Color.enumMethod", 24, 25, false},
		"area":       {"com.example.app.Shape.area", 43, 43, true},
		"Point":      {"com.example.app.Point.Point", 47, 47, false},
	}

	if len(funcs) != len(wants) {
		t.Fatalf("got %d funcs, want %d: %+v", len(funcs), len(wants), funcs)
	}
	for _, f := range funcs {
		w, ok := wants[f.Name]
		if !ok {
			t.Errorf("unexpected func %q (nameless/empty/field must be skipped)", f.Name)
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

// TestExtractJavaFunctionsNilRoot covers the nil-root guard.
func TestExtractJavaFunctionsNilRoot(t *testing.T) {
	if funcs := extractJavaFunctions(nil, "x.java", ""); funcs != nil {
		t.Errorf("nil root: got %+v, want nil", funcs)
	}
}

// TestJavaPackageNameDefault covers the default-package (no declaration) branch.
func TestJavaPackageNameDefault(t *testing.T) {
	const xml = `<?xml version="1.0"?>
<sources>
  <source name="/p/D.java">
    <program srow="0" scol="0" erow="3" ecol="0">
      <class_declaration srow="0" scol="0" erow="2" ecol="1">
        <identifier field="name" srow="0" scol="6" erow="0" ecol="7">D</identifier>
        <class_body field="body" srow="0" scol="8" erow="2" ecol="1">
          <method_declaration srow="1" scol="2" erow="1" ecol="20">
            <identifier field="name" srow="1" scol="2" erow="1" ecol="3">m</identifier>
          </method_declaration>
        </class_body>
      </class_declaration>
    </program>
  </source>
</sources>`
	sources, err := treesitter.ParseXML([]byte(xml))
	if err != nil {
		t.Fatalf("ParseXML: %v", err)
	}
	if pkg := javaPackageName(sources[0].Root); pkg != "" {
		t.Errorf("javaPackageName = %q, want empty", pkg)
	}
	funcs := extractJavaFunctions(sources[0].Root, "D.java", "")
	if len(funcs) != 1 || funcs[0].QualifiedName != "D.m" {
		t.Fatalf("default-package funcs = %+v, want single D.m", funcs)
	}
}
