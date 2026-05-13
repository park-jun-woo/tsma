//ff:func feature=endpoint type=implementation control=sequence
//ff:what Builds a tsRouteRegistration from a regex match on Express route pattern
package endpoint

func buildExpressRegistration(content string, lines []string, loc []int, relPath string) tsRouteRegistration {
	method := content[loc[2]:loc[3]]
	routePath := content[loc[4]:loc[5]]
	handlerName := content[loc[6]:loc[7]]

	regLine := countNewlines(content[:loc[0]]) + 1

	startLine, endLine := findTSFuncLocation(lines, handlerName)
	if startLine == 0 {
		startLine = regLine
		endLine = regLine
	}

	return tsRouteRegistration{
		method:    method,
		path:      routePath,
		handler:   handlerName,
		file:      relPath,
		startLine: startLine,
		endLine:   endLine,
	}
}
