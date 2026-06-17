//ff:func feature=gate type=helper control=iteration dimension=1 lang=typescript
//ff:what scanTSSmells: the TypeScript analogue of scanGoSmells. Joins each matched root-relative test file with root and runs smell.ScanTS (tree-sitter node-based). Errors (tree-sitter absent or unparseable file) are ignored exactly as the Go path ignores parse errors — TS smell is LevelReview/best-effort.
package tsmagate

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/smell"
)

// scanTSSmells statically scans the matched TS/JS test files for escape-hatch
// smells and returns the collected findings. A scan error is ignored (continue):
// a broken or unparseable test file is judged by tests-must-pass, not here.
func scanTSSmells(root string, files []string) []smell.Finding {
	var findings []smell.Finding
	for _, file := range files {
		found, err := smell.ScanTS(filepath.Join(root, file))
		if err != nil {
			continue
		}
		findings = append(findings, found...)
	}
	return findings
}
