//ff:func feature=index type=helper control=sequence lang=typescript
//ff:what tsArrowFuncFromDeclarator: builds a model.Function from a variable_declarator whose value is an arrow_function/function_expression (`const foo = () => ...`), or returns false. The declarator span is the range, so a multi-line arrow body is captured accurately.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// tsArrowFuncFromDeclarator converts a function-valued declarator into a
// model.Function, or returns false when it is not a function binding.
func tsArrowFuncFromDeclarator(d *treesitter.Node, relDir, relPath string, exported bool) (model.Function, bool) {
	if d.Type != "variable_declarator" {
		return model.Function{}, false
	}
	val := d.ChildByField("value")
	if val == nil {
		return model.Function{}, false
	}
	if val.Type != "arrow_function" && val.Type != "function_expression" && val.Type != "function" {
		return model.Function{}, false
	}
	nameNode := d.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return model.Function{}, false
	}
	name := nameNode.Text
	return model.Function{
		QualifiedName: buildQualifiedName(relDir, "", name),
		Name:          name,
		File:          relPath,
		StartLine:     d.StartLine(),
		EndLine:       d.EndLine(),
		Exported:      exported,
		Status:        model.StatusTodo,
	}, true
}
