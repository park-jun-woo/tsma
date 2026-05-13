//ff:func feature=endpoint type=implementation control=sequence
//ff:what Collects FastAPI route registrations and converts them to endpoints
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

// Detect scans all .py files for FastAPI route registrations.
func (d *PyFastapiDetector) Detect(projectRoot string) ([]model.Endpoint, error) {
	routes, err := collectFastapiRoutes(projectRoot)
	if err != nil {
		return nil, err
	}

	endpoints := make([]model.Endpoint, 0, len(routes))
	for _, r := range routes {
		endpoints = append(endpoints, pyRouteToEndpoint(r))
	}
	return endpoints, nil
}
