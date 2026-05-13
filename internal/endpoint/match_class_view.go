//ff:func feature=endpoint type=implementation control=sequence
//ff:what Matches a Django class-based view definition and returns endpoints per HTTP method
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func matchClassView(line string, lines []string, idx, lineNum int, relPath string, route djangoRoute) []model.Endpoint {
	m := djangoClassViewRe.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	if m[2] != route.viewName {
		return nil
	}

	classIndent := m[1]
	methods := findClassMethods(lines, idx, classIndent)
	if len(methods) == 0 {
		return buildAnyEndpoint(lines, idx, lineNum, classIndent, relPath, route)
	}

	return buildClassMethodEndpoints(methods, relPath, route)
}
