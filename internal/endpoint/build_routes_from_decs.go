//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Builds pyRoute entries from pending decorators matched to a def line
package endpoint

func buildRoutesFromDecs(pendingDec []pendingDecorator, dm []string, lines []string, idx, lineNum int, relPath string) []pyRoute {
	funcIndent := dm[1]
	funcName := dm[2]
	endLine := findPyFuncEnd(lines, idx, funcIndent)

	var routes []pyRoute
	for _, pd := range pendingDec {
		routes = append(routes, pyRoute{
			method:    pd.method,
			path:      pd.path,
			handler:   funcName,
			file:      relPath,
			startLine: lineNum,
			endLine:   endLine,
		})
	}
	return routes
}
