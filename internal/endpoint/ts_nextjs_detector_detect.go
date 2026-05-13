//ff:func feature=endpoint type=implementation control=sequence
//ff:what Scans Next.js project for App Router and Pages Router API endpoints
package endpoint

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Detect scans a Next.js project for API route endpoints.
func (d *TSNextjsDetector) Detect(projectRoot string) ([]model.Endpoint, error) {
	var endpoints []model.Endpoint

	appRouterEps, err := detectAppRouterRoutes(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("detect app router routes: %w", err)
	}
	endpoints = append(endpoints, appRouterEps...)

	pagesRouterEps, err := detectPagesRouterRoutes(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("detect pages router routes: %w", err)
	}
	endpoints = append(endpoints, pagesRouterEps...)

	return endpoints, nil
}
