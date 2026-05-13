//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Finds a Python function by file path and start line in the function index
package chain

// findPyFuncAtLocation finds a function at the given file and start line.
func findPyFuncAtLocation(funcs map[string]*pyFuncInfo, file string, startLine int) *pyFuncInfo {
	for _, fi := range funcs {
		if fi.file == file && fi.startLine == startLine {
			return fi
		}
	}
	// Fallback: match by file only if there's exactly one function.
	var matched *pyFuncInfo
	count := 0
	for _, fi := range funcs {
		if fi.file == file {
			matched = fi
			count++
		}
	}
	if count == 1 {
		return matched
	}
	return nil
}
