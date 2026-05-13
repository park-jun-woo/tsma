//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Searches for a Python callee matching class.method pattern in the index
package chain

import "strings"

// findPyCalleeByClass searches for a callee matching class.method pattern.
func findPyCalleeByClass(funcs map[string]*pyFuncInfo, funcName, className string) *pyFuncInfo {
	for key, fi := range funcs {
		if fi.name == funcName && strings.Contains(key, className+".") {
			return fi
		}
	}
	return nil
}
