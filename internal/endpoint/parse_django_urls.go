//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Parses a single urls.py file for path() and re_path() patterns
package endpoint

import "os"

// parseDjangoURLs parses a single urls.py file for path() and re_path() patterns.
func parseDjangoURLs(filePath string) ([]djangoRoute, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	content := string(data)

	var routes []djangoRoute

	for _, m := range djangoPathRe.FindAllStringSubmatch(content, -1) {
		routes = append(routes, djangoRoute{
			path:     m[1],
			viewName: extractViewName(m[2]),
		})
	}

	for _, m := range djangoRePathRe.FindAllStringSubmatch(content, -1) {
		routes = append(routes, djangoRoute{
			path:     m[1],
			viewName: extractViewName(m[2]),
		})
	}

	return routes, nil
}
