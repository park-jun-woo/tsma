//ff:func feature=endpoint type=implementation control=sequence
//ff:what Collects Express route registrations and converts them to endpoints
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

// Detect scans all TypeScript and JavaScript source files for Express route registrations.
func (d *TSExpressDetector) Detect(projectRoot string) ([]model.Endpoint, error) {
	regs, err := collectExpressRoutes(projectRoot)
	if err != nil {
		return nil, err
	}

	endpoints := make([]model.Endpoint, 0, len(regs))
	for _, r := range regs {
		endpoints = append(endpoints, tsRouteToEndpoint(r))
	}
	return endpoints, nil
}
