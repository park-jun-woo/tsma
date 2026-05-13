//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Finds a Go function by file path and start line in the function index
package chain

// findFuncAtLocation finds a function at the given file and line.
func findFuncAtLocation(funcs map[string]*funcInfo, file string, startLine int) *funcInfo {
	for _, fi := range funcs {
		if fi.file == file && fi.startLine == startLine {
			return fi
		}
	}
	// Fallback: match by file only if there's exactly one function.
	var matched *funcInfo
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
