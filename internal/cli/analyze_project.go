//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Detects language, indexes functions, and matches test files to build a session
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/park-jun-woo/tsma/internal/detect"
	"github.com/park-jun-woo/tsma/internal/index"
	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzeProject detects the language, indexes functions, and matches test files.
func analyzeProject(projectRoot string) (*model.Session, error) {
	lf, err := detect.Detect(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("detect language: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Detected: %s\n", lf.Lang)

	fmt.Fprintln(os.Stderr, "Indexing functions...")
	idxr := index.NewIndexer(lf.Lang)
	functions, err := idxr.Index(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("index functions: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Found %d functions\n", len(functions))

	fmt.Fprintln(os.Stderr, "Matching test files...")
	for i := range functions {
		functions[i].Status = model.StatusTodo
	}
	attributeTests(projectRoot, lf.Lang, functions)

	sess := &model.Session{
		Project:   projectRoot,
		Lang:      lf.Lang,
		CheckedAt: time.Now(),
		Functions: functions,
	}

	// First-scan batch: measure every function's coverage now so a single
	// `tsma next` yields PASS for 100% functions and measured TODOs for the
	// rest. After this the session is fully measured, so first-pass mode is done
	// and `tsma next` goes straight to the interactive rotating cursor.
	fmt.Fprintln(os.Stderr, "Measuring coverage...")
	batchMeasure(projectRoot, sess)
	sess.FirstPassDone = true
	sess.CurrentIndex = 0
	sess.RecalcSummary()

	return sess, nil
}
