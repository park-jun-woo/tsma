//ff:func feature=endpoint type=implementation control=sequence
//ff:what Collects Chi route registrations and converts them to endpoints
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

// Detect scans all Go source files for Chi router registrations.
func (d *GoChiDetector) Detect(projectRoot string) ([]model.Endpoint, error) {
	regs, err := collectChiRoutes(projectRoot)
	if err != nil {
		return nil, err
	}

	endpoints := make([]model.Endpoint, 0, len(regs))
	for _, r := range regs {
		endpoints = append(endpoints, routeToEndpoint(r))
	}
	return endpoints, nil
}
