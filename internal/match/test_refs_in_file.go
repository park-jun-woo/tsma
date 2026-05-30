//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Maps each TestXxx in a parsed file to the source identifiers it references
package match

import "go/ast"

// testRefsInFile maps each top-level TestXxx function in funcs to the set of
// source identifier references (name + statically-resolved receiver type) it
// makes. References include identifiers called directly in the test body plus,
// via a same-file 1-hop call graph, the identifiers called by any non-Test
// helper that the test invokes directly.
func testRefsInFile(funcs map[string]*ast.FuncDecl) map[string]map[calledRef]struct{} {
	result := make(map[string]map[calledRef]struct{})
	for name, fd := range funcs {
		if !isTestFuncName(name) {
			continue
		}
		result[name] = expandTestRefs(fd, funcs)
	}
	return result
}
