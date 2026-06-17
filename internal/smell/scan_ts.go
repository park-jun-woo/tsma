//ff:func feature=smell type=helper control=sequence level=review lang=typescript
//ff:what ScanTS: parses one TS/JS test file with tree-sitter (node-based, never substring) and runs the three TS escape-hatch detectors (as-any private bypass, Reflect., Object.getOwnProperty*). Returns (nil, err) when tree-sitter is unavailable or the file fails to parse, so the caller (scanTSSmells) ignores it exactly as it ignores a Go parse error — smell is LevelReview/best-effort. Mirrors ScanGo.
package smell

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// ScanTS parses a single TS/JS test file with tree-sitter and collects the TS
// escape-hatch findings. It returns (nil, err) when the precise path cannot run
// so the caller can ignore it (no substring fallback — that would re-introduce
// the false positives the node-based approach eliminates).
func ScanTS(path string) ([]Finding, error) {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil, fmt.Errorf("smell: tree-sitter unavailable for %s", path)
	}
	grammar := treesitter.ResolveGrammar("typescript")
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	root, err := treesitter.ParseFile(command, grammar, abs)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	findings = append(findings, detectTSAsAny(root, path)...)
	findings = append(findings, detectTSReflect(root, path)...)
	findings = append(findings, detectTSOwnProperty(root, path)...)
	return findings, nil
}
