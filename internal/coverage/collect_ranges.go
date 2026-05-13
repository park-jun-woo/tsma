//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Extracts Go function ranges from the endpoint handler and chain entries
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

func collectRanges(ep *model.Endpoint) []funcRange {
	var ranges []funcRange

	if ep.Handler.File != "" {
		ranges = append(ranges, funcRange{
			file:      ep.Handler.File,
			startLine: ep.Handler.StartLine,
			endLine:   ep.Handler.EndLine,
			funcName:  ep.Name,
		})
	}

	for _, ce := range ep.Chain {
		if ce.File != "" && ce.Boundary == "" {
			ranges = append(ranges, funcRange{
				file:      ce.File,
				startLine: ce.StartLine,
				endLine:   ce.EndLine,
				funcName:  ce.Func,
			})
		}
	}

	return ranges
}
