//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Appends pending decorators from an api_route regex match
package endpoint

func appendAPIRouteDecs(pending []pendingDecorator, m []string, lineNum int) []pendingDecorator {
	methods := parseMethodsList(m[2])
	for _, method := range methods {
		pending = append(pending, pendingDecorator{
			method: method,
			path:   m[1],
			line:   lineNum,
		})
	}
	return pending
}
