//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Performs initial project analysis detecting language, endpoints, and call chains
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/park-jun-woo/tsma/internal/chain"
	"github.com/park-jun-woo/tsma/internal/detect"
	"github.com/park-jun-woo/tsma/internal/endpoint"
	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzeProject performs initial project analysis.
func analyzeProject(projectRoot string) (*model.Session, error) {
	// Step 1: Detect language and framework.
	lf, err := detect.Detect(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("detect language: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Detected: %s (%s)\n", lf.Lang, lf.Framework)

	// Step 2: Detect endpoints.
	det := endpoint.NewDetector(lf.Lang, lf.Framework)
	endpoints, err := det.Detect(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("detect endpoints: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Found %d endpoints\n", len(endpoints))

	// Step 3: Trace call chains for each endpoint.
	tracer := chain.NewTracer(lf.Lang)
	for i := range endpoints {
		ep := &endpoints[i]
		if ep.Handler.File == "" {
			continue
		}
		entries, err := tracer.Trace(projectRoot, ep.Handler)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: chain trace failed for %s: %v\n", ep.Name, err)
			continue
		}
		ep.Chain = entries
	}

	sess := &model.Session{
		Project:   projectRoot,
		Lang:      lf.Lang,
		Framework: lf.Framework,
		Created:   time.Now(),
		Endpoints: endpoints,
	}
	sess.RecalcSummary()

	return sess, nil
}
