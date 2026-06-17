//ff:func feature=match type=helper control=iteration dimension=1 lang=rust
//ff:what ingestRsTestsDir: scans the project's tests/ directory for *.rs integration tests and records, for every called name in each (the whole file is test code), a back-reference to that test file in the index — the tests/ fallback attribution for pub functions. A missing tests/ dir is a no-op (graceful — the index then carries only in-file references). Parse failures per file are skipped.
package match

import (
	"os"
	"path/filepath"
	"strings"
)

// ingestRsTestsDir ingests every tests/*.rs integration test into idx.
func ingestRsTestsDir(idx *RsTestIndex, projectRoot, command, grammar string) {
	dir := filepath.Join(projectRoot, "tests")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".rs") {
			continue
		}
		ingestRsTestFile(idx, command, grammar, filepath.Join(dir, e.Name()), filepath.Join("tests", e.Name()))
	}
}
