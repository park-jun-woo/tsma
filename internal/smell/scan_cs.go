//ff:func feature=smell type=helper control=sequence level=review lang=csharp
//ff:what ScanCs: parses one C# test file with tree-sitter (node-based, never substring) and runs the two C# escape-hatch detectors (System.Reflection GetMethod/GetField/GetProperty, MethodInfo/PropertyInfo/FieldInfo dynamic-invocation declarations). Returns (nil, err) when tree-sitter is unavailable or the file fails to parse, so the caller (scanCsSmells) ignores it exactly as it ignores a Go parse error — smell is LevelReview/best-effort. Mirrors ScanJava. AF015 original .NET patterns.
package smell

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// ScanCs parses a single C# test file with tree-sitter and collects the C#
// escape-hatch findings. It returns (nil, err) when the precise path cannot run
// so the caller can ignore it (no substring fallback — that would re-introduce
// the false positives the node-based approach eliminates).
func ScanCs(path string) ([]Finding, error) {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil, fmt.Errorf("smell: tree-sitter unavailable for %s", path)
	}
	grammar := treesitter.ResolveGrammar("csharp")
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	root, err := treesitter.ParseFile(command, grammar, abs)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	findings = append(findings, detectCsReflect(root, path)...)
	findings = append(findings, detectCsReflectInfo(root, path)...)
	return findings, nil
}
