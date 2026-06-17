//ff:func feature=match type=helper control=iteration dimension=1 lang=typescript
//ff:what ingestTSDir: scans one directory for *.test/*.spec files and ingests each into the content index (via ingestTSTestFile). Extracted from BuildTSPkgTestIndex so the dir-loop and the entry-loop stay at separate depths (Q1 ≤2). A missing directory is a no-op.
package match

import (
	"os"
	"path/filepath"
)

// ingestTSDir ingests every TS/JS test file directly in dir (relative to
// projectRoot) into idx.
func ingestTSDir(idx *TSPkgTestIndex, projectRoot, dir, command, grammar string) {
	absDir := filepath.Join(projectRoot, dir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !isTSTestFile(name) {
			continue
		}
		abs := filepath.Join(absDir, name)
		ingestTSTestFile(idx, command, grammar, abs, relTestPath(projectRoot, abs))
	}
}
