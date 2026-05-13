//ff:func feature=cli type=helper control=sequence
//ff:what Performs initial project analysis: detect language and index functions
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/park-jun-woo/tsma/internal/detect"
	"github.com/park-jun-woo/tsma/internal/index"
	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzeProject performs initial project analysis.
func analyzeProject(projectRoot string) (*model.Session, error) {
	// Step 1: Detect language.
	lf, err := detect.Detect(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("detect language: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Detected: %s\n", lf.Lang)

	// Step 2: Index all functions.
	fmt.Fprintln(os.Stderr, "Indexing functions...")
	idxr := index.NewIndexer(lf.Lang)
	functions, err := idxr.Index(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("index functions: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Found %d functions\n", len(functions))

	// Step 3: Set initial status for all functions.
	setInitialStatus(functions)

	sess := &model.Session{
		Project:   projectRoot,
		Lang:      lf.Lang,
		Created:   time.Now(),
		Functions: functions,
	}
	sess.RecalcSummary()

	return sess, nil
}
