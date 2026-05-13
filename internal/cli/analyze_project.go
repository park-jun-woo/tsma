//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Performs initial project analysis: detect language, index functions, build call graph
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/park-jun-woo/tsma/internal/detect"
	"github.com/park-jun-woo/tsma/internal/graph"
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

	// Step 3: Build call graph.
	fmt.Fprintln(os.Stderr, "Building call graph...")
	bldr := graph.NewBuilder(lf.Lang)
	functions, gs, err := bldr.Build(projectRoot, functions)
	if err != nil {
		return nil, fmt.Errorf("build call graph: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Graph: %d nodes, %d edges, %d entry points, %d dead\n",
		gs.Nodes, gs.Edges, gs.EntryPoints, gs.Dead)

	// Step 4: Set initial status for all functions.
	for i := range functions {
		if functions[i].Status == "" {
			functions[i].Status = model.StatusTodo
		}
	}

	sess := &model.Session{
		Project:   projectRoot,
		Lang:      lf.Lang,
		Created:   time.Now(),
		Functions: functions,
		Graph:     gs,
	}
	sess.RecalcSummary()

	return sess, nil
}
