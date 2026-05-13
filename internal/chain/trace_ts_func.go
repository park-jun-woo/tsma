//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Extracts calls from a TS/JS function body and recurses into found definitions
package chain

import "github.com/park-jun-woo/tsma/internal/model"

// traceTSFunc extracts function calls from the given line range and recurses into found definitions.
func traceTSFunc(projectRoot, file string, lines []string, startLine, endLine int, visited map[string]bool, entries *[]model.ChainEntry, depth int) {
	if depth >= maxTraceDepth {
		return
	}

	body := extractTSBody(lines, startLine, endLine)
	if body == "" {
		return
	}

	matches := tsFuncCallPattern.FindAllStringSubmatch(body, -1)
	for _, m := range matches {
		processTSCallMatch(m, projectRoot, file, visited, entries, depth)
	}
}
