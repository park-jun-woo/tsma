//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Searches the Python function index for a matching callee by name or class prefix
package chain

// findPyCallee searches the function index for a matching callee.
func findPyCallee(funcs map[string]*pyFuncInfo, funcName string, parts []string) *pyFuncInfo {
	// Direct function name match.
	for _, fi := range funcs {
		if fi.name == funcName {
			return fi
		}
	}

	// Try matching with class prefix: ClassName.method
	if len(parts) >= 2 {
		return findPyCalleeByClass(funcs, funcName, parts[len(parts)-2])
	}

	return nil
}
