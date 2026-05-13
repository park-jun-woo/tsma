//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Builds endpoints from class-based view HTTP method definitions
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func buildClassMethodEndpoints(methods []classMethod, relPath string, route djangoRoute) []model.Endpoint {
	var endpoints []model.Endpoint
	for _, cm := range methods {
		endpoints = append(endpoints, classMethodToEndpoint(cm, relPath, route))
	}
	return endpoints
}
