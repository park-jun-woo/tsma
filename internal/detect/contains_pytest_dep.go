//ff:func feature=detect type=helper control=iteration dimension=1 lang=python
//ff:what Reports whether pyproject.toml declares pytest as a PEP 621 dependency
package detect

import (
	"os"
	"path/filepath"
	"strings"
)

// containsPytestDep reports whether the project's pyproject.toml declares pytest
// in a PEP 621 dependency block: [project.optional-dependencies],
// [dependency-groups], or [project]'s dependencies array (inline or multiline).
//
// Detection is a deliberately simple line scan rather than a full TOML parse
// (over-engineering avoided, per Phase006 §2/§4): once a recognised dependency
// context opens, any subsequent line mentioning "pytest" — until that context
// ends — counts. This matches the existing containsPytest token-scan fidelity;
// edge cases like commented-out pytest are accepted as out of scope. Per-line
// state transitions live in advanceDepScan to keep nesting at depth ≤2 (Q1).
func containsPytestDep(projectRoot string) bool {
	data, err := os.ReadFile(filepath.Join(projectRoot, "pyproject.toml"))
	if err != nil {
		return false
	}

	var st depScanState
	for _, raw := range strings.Split(string(data), "\n") {
		hit, next := advanceDepScan(st, strings.TrimSpace(raw))
		if hit {
			return true
		}
		st = next
	}
	return false
}
