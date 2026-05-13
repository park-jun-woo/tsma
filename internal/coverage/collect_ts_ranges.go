//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Collects handler and chain function file/line ranges for TS coverage analysis
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectTSRanges extracts function ranges from an endpoint's handler and chain.
func collectTSRanges(ep *model.Endpoint) []tsFuncRange {
	var ranges []tsFuncRange

	if ep.Handler.File != "" {
		ranges = append(ranges, tsFuncRange{
			file:      ep.Handler.File,
			startLine: ep.Handler.StartLine,
			endLine:   ep.Handler.EndLine,
			funcName:  ep.Name,
		})
	}

	for _, ce := range ep.Chain {
		if ce.File != "" && ce.Boundary == "" {
			ranges = append(ranges, tsFuncRange{
				file:      ce.File,
				startLine: ce.StartLine,
				endLine:   ce.EndLine,
				funcName:  ce.Func,
			})
		}
	}

	return ranges
}
