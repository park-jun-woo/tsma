//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Calculates average coverage percentage across PASS and DONE functions
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// computeAverageCoverage calculates average coverage across PASS and DONE functions.
func computeAverageCoverage(sess *model.Session) float64 {
	var sum float64
	var count int
	for _, fn := range sess.Functions {
		if fn.Status == model.StatusPass || fn.Status == model.StatusDone {
			sum += fn.CoveragePct
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}
