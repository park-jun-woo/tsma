//ff:func feature=gate type=helper control=selection
//ff:what scanSmells: dispatches escape-hatch smell scanning to the language detector — Go via go/ast (scanGoSmells), TypeScript via tree-sitter (scanTSSmells). Returns nil for languages without a detector. Collapses Prepare's per-language smell branch into one call (keeps Prepare control=sequence).
package tsmagate

import "github.com/park-jun-woo/tsma/internal/smell"

// scanSmells statically scans the matched test files for the language's
// escape-hatch smells, or returns nil when the language has no detector.
func scanSmells(lang, root string, files []string) []smell.Finding {
	switch lang {
	case "go":
		return scanGoSmells(root, files)
	case "typescript":
		return scanTSSmells(root, files)
	}
	return nil
}
