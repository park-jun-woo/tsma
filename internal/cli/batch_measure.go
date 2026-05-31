//ff:func feature=cli type=helper control=sequence
//ff:what Batch-measures coverage for every function at analysis time, dispatching Go vs non-Go
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// batchMeasure measures the coverage of every indexed function in one pass at
// analysis time (the first-scan batch). 100% functions become PASS; functions
// with a matched test below 100% become measured TODOs (Attempt=1); untested
// functions stay TODO. It never promotes a partial to DONE — auto-DONE
// convergence is left to later interactive re-measurement. Go uses a
// package-at-a-time run (one `go test` per package); other languages fall back
// to a per-function Check loop. After it returns the session is fully measured.
func batchMeasure(root string, sess *model.Session) {
	if sess.Lang == "go" {
		batchMeasureGo(root, sess)
		return
	}
	batchMeasureOther(root, sess)
}
