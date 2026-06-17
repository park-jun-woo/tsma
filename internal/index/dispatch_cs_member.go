//ff:func feature=index type=helper control=selection dimension=1 lang=csharp
//ff:what dispatchCSMember: classifies one member node of a C# compilation unit / namespace / type body — method_declaration / constructor_declaration / destructor_declaration / property_declaration are emitted as model.Functions scoped by the current namespace+type stack (csFuncFromNode); block namespace_declaration and nested class/struct/interface/record/enum declarations recurse via walkCSTypeDecl (pushing a scope from the dotted name). A file_scoped_namespace_declaration carries no members (handled as fileNs by extractCSharpFunctions) and is ignored here. The single place C#'s indexable shapes are recognized (the dispatchJavaMember analogue).
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// dispatchCSMember handles one member node by type.
func dispatchCSMember(c *treesitter.Node, fileNs string, scopes []csScope, relPath string, out *[]model.Function) {
	switch c.Type {
	case "method_declaration", "constructor_declaration", "destructor_declaration",
		"property_declaration":
		if fn, ok := csFuncFromNode(c, fileNs, scopes, relPath); ok {
			*out = append(*out, fn)
		}
	case "namespace_declaration", "class_declaration", "struct_declaration",
		"interface_declaration", "record_declaration", "record_struct_declaration",
		"enum_declaration":
		walkCSTypeDecl(c, fileNs, scopes, relPath, out)
	}
}
