//ff:func feature=index type=helper control=sequence
//ff:what Attempts to match a TS class method when inside a class context
package index

import "github.com/park-jun-woo/tsma/internal/model"

// tryMatchTSMethod attempts to match a class method when inside a class context.
func tryMatchTSMethod(line, currentClass, relDir, relPath string, lineNum int) (model.Function, bool) {
	if currentClass == "" {
		return model.Function{}, false
	}
	return matchTSMethod(line, currentClass, relDir, relPath, lineNum)
}
