//ff:func feature=smell type=helper control=iteration dimension=1 level=review lang=rust
//ff:what ScanRs: parses one Rust file with tree-sitter (node-based, never substring), restricts attention to its test scopes (rsTestScopeNodes — in-file #[cfg(test)] modules + top-level #[test] fns), and runs the three Rust escape-hatch detectors (unsafe, transmute, std::ptr) over each. Scoping to test code is what keeps legitimate production/FFI unsafe from ever producing a finding (false-positive zero). Returns (nil, err) when tree-sitter is unavailable or the file fails to parse, so the caller (scanRsSmells) ignores it exactly as the Go path ignores parse errors — smell is LevelReview/best-effort. Mirrors ScanCs.
package smell

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// ScanRs parses a single Rust file with tree-sitter and collects the Rust
// escape-hatch findings from its test scopes. It returns (nil, err) when the
// precise path cannot run so the caller can ignore it (no substring fallback —
// that would re-introduce the false positives the node-based approach eliminates).
func ScanRs(path string) ([]Finding, error) {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil, fmt.Errorf("smell: tree-sitter unavailable for %s", path)
	}
	grammar := treesitter.ResolveGrammar("rust")
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	root, err := treesitter.ParseFile(command, grammar, abs)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	for _, scope := range rsTestScopeNodes(root) {
		findings = append(findings, detectRsUnsafe(scope, path)...)
		findings = append(findings, detectRsTransmute(scope, path)...)
		findings = append(findings, detectRsPtr(scope, path)...)
	}
	return findings, nil
}
