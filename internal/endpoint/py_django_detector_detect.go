//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Walks Django project to find urls.py files and resolves view definitions
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

// Detect scans Django projects for urlpatterns and resolves view functions.
func (d *PyDjangoDetector) Detect(projectRoot string) ([]model.Endpoint, error) {
	urlsFiles, err := collectDjangoURLFiles(projectRoot)
	if err != nil {
		return nil, err
	}

	var urlRoutes []djangoRoute
	for _, urlsFile := range urlsFiles {
		routes, parseErr := parseDjangoURLs(urlsFile)
		if parseErr != nil {
			continue
		}
		urlRoutes = append(urlRoutes, routes...)
	}

	var endpoints []model.Endpoint
	for _, route := range urlRoutes {
		eps := resolveDjangoView(projectRoot, route)
		endpoints = append(endpoints, eps...)
	}

	return endpoints, nil
}
