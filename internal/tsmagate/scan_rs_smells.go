//ff:func feature=gate type=helper control=iteration dimension=1 lang=rust
//ff:what scanRsSmells: the Rust analogue of scanCsSmells/scanGoSmells. Joins each matched root-relative test file with root and runs smell.ScanRs (tree-sitter node-based, scoped to test code). Errors (tree-sitter absent or unparseable file) are ignored exactly as the Go path ignores parse errors — Rust smell is LevelReview/best-effort. For in-file Rust tests the "test file" is the source file itself, which ScanRs restricts to its #[cfg(test)] scopes.
package tsmagate

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/smell"
)

// scanRsSmells statically scans the matched Rust test files for escape-hatch
// smells and returns the collected findings. A scan error is ignored (continue):
// a broken or unparseable file is judged by tests-must-pass, not here.
func scanRsSmells(root string, files []string) []smell.Finding {
	var findings []smell.Finding
	for _, file := range files {
		found, err := smell.ScanRs(filepath.Join(root, file))
		if err != nil {
			continue
		}
		findings = append(findings, found...)
	}
	return findings
}
