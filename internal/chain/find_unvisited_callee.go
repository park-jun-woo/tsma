//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Searches the Go function index for a matching unvisited callee
package chain

// findUnvisitedCallee searches the index for a matching unvisited callee.
func findUnvisitedCallee(funcs map[string]*funcInfo, calleeName string, visited map[string]bool) *funcInfo {
	for key, fi := range funcs {
		if fi.name == calleeName && !visited[key] {
			return fi
		}
	}
	return nil
}
