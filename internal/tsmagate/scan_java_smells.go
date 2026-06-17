//ff:func feature=gate type=helper control=iteration dimension=1 lang=java
//ff:what scanJavaSmells: the Java analogue of scanGoSmells/scanTSSmells. Joins each matched root-relative test file with root and runs smell.ScanJava (tree-sitter node-based). Errors (tree-sitter absent or unparseable file) are ignored exactly as the Go path ignores parse errors — Java smell is LevelReview/best-effort.
package tsmagate

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/smell"
)

// scanJavaSmells statically scans the matched Java test files for escape-hatch
// smells and returns the collected findings. A scan error is ignored (continue):
// a broken or unparseable test file is judged by tests-must-pass, not here.
func scanJavaSmells(root string, files []string) []smell.Finding {
	var findings []smell.Finding
	for _, file := range files {
		found, err := smell.ScanJava(filepath.Join(root, file))
		if err != nil {
			continue
		}
		findings = append(findings, found...)
	}
	return findings
}
