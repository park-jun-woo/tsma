//ff:func feature=detect type=helper control=iteration dimension=1 lang=python
//ff:what Reports whether a tests/ dir holds pytest-style test_*.py / *_test.py files
package detect

import (
	"os"
	"path/filepath"
	"strings"
)

// hasPytestLayout reports whether the project follows the conventional pytest
// layout: a tests/ directory containing at least one test_*.py or *_test.py
// file. Only the immediate tests/ directory is scanned (no recursion) — a
// cheap, high-signal heuristic.
func hasPytestLayout(projectRoot string) bool {
	entries, err := os.ReadDir(filepath.Join(projectRoot, "tests"))
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".py") {
			continue
		}
		if strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test.py") {
			return true
		}
	}
	return false
}
