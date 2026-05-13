//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Extracts Python function ranges from the endpoint handler and chain entries
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectPyRanges extracts function ranges from the endpoint.
func collectPyRanges(ep *model.Endpoint) []pyFuncRange {
	var ranges []pyFuncRange

	if ep.Handler.File != "" {
		ranges = append(ranges, pyFuncRange{
			file:      ep.Handler.File,
			startLine: ep.Handler.StartLine,
			endLine:   ep.Handler.EndLine,
			funcName:  ep.Name,
		})
	}

	for _, ce := range ep.Chain {
		if ce.File != "" && ce.Boundary == "" {
			ranges = append(ranges, pyFuncRange{
				file:      ce.File,
				startLine: ce.StartLine,
				endLine:   ce.EndLine,
				funcName:  ce.Func,
			})
		}
	}

	return ranges
}
