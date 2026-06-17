//ff:func feature=smell type=helper control=sequence level=review lang=java
//ff:what ScanJava: parses one JUnit test file with tree-sitter (node-based, never substring) and runs the two Java escape-hatch detectors (java.lang.reflect getDeclaredMethod/Field, setAccessible(true)). Returns (nil, err) when tree-sitter is unavailable or the file fails to parse, so the caller (scanJavaSmells) ignores it exactly as it ignores a Go parse error — smell is LevelReview/best-effort. Mirrors ScanTS.
package smell

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// ScanJava parses a single Java test file with tree-sitter and collects the Java
// escape-hatch findings. It returns (nil, err) when the precise path cannot run
// so the caller can ignore it (no substring fallback — that would re-introduce
// the false positives the node-based approach eliminates).
func ScanJava(path string) ([]Finding, error) {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil, fmt.Errorf("smell: tree-sitter unavailable for %s", path)
	}
	grammar := treesitter.ResolveGrammar("java")
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	root, err := treesitter.ParseFile(command, grammar, abs)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	findings = append(findings, detectJavaReflect(root, path)...)
	findings = append(findings, detectJavaSetAccessible(root, path)...)
	return findings, nil
}
