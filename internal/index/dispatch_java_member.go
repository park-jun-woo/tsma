//ff:func feature=index type=helper control=selection dimension=1 lang=java
//ff:what dispatchJavaMember: classifies one member node of a Java type body — method_declaration / constructor_declaration are emitted as model.Functions scoped by the current type stack (javaFuncFromNode); nested class/interface/enum/record/annotation declarations recurse via walkJavaTypeDecl (pushing a scope); an enum_body_declarations wrapper is descended without a new scope. The single place Java's indexable shapes are recognized (the walkTSChild analogue).
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// dispatchJavaMember handles one member node by type.
func dispatchJavaMember(c *treesitter.Node, pkg string, scopes []javaScope, relPath string, out *[]model.Function) {
	switch c.Type {
	case "method_declaration", "constructor_declaration", "compact_constructor_declaration":
		if fn, ok := javaFuncFromNode(c, pkg, scopes, relPath); ok {
			*out = append(*out, fn)
		}
	case "class_declaration", "interface_declaration", "enum_declaration",
		"record_declaration", "annotation_type_declaration":
		walkJavaTypeDecl(c, pkg, scopes, relPath, out)
	case "enum_body_declarations":
		collectJavaMembers(c, pkg, scopes, relPath, out)
	}
}
