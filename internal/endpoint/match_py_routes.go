//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Matches Python decorator and def lines to build FastAPI route entries
package endpoint

import "strings"

func matchPyRoutes(lines []string, relPath string) []pyRoute {
	var (
		routes     []pyRoute
		pendingDec []pendingDecorator
	)

	for i, line := range lines {
		lineNum := i + 1

		if m := pyDecoratorRe.FindStringSubmatch(line); m != nil {
			pendingDec = append(pendingDec, pendingDecorator{
				method: strings.ToUpper(m[1]),
				path:   m[2],
				line:   lineNum,
			})
			continue
		}

		if m := pyAPIRouteRe.FindStringSubmatch(line); m != nil {
			pendingDec = appendAPIRouteDecs(pendingDec, m, lineNum)
			continue
		}

		newRoutes, clearedDec := resolvePendingDecs(pendingDec, lines, line, i, lineNum, relPath)
		routes = append(routes, newRoutes...)
		pendingDec = clearedDec
	}

	return routes
}
