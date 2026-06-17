//ff:func feature=match type=helper control=iteration dimension=1 lang=csharp
//ff:what ingestCsTestDir: scans one test directory for *Tests.cs/*Test.cs files and ingests each into the content index (via ingestCsTestFile). Extracted from BuildCsPkgTestIndex so the dir scan and the entry loop stay at separate depths (Q1 ≤2). A missing directory is a no-op (graceful fallback to filename matching upstream).
package match

import (
	"os"
	"path/filepath"
)

// ingestCsTestDir ingests every C# test file directly in dir (relative to
// projectRoot) into idx.
func ingestCsTestDir(idx *CsPkgTestIndex, projectRoot, dir, command, grammar string) {
	absDir := filepath.Join(projectRoot, dir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !isCsTestFile(name) {
			continue
		}
		abs := filepath.Join(absDir, name)
		ingestCsTestFile(idx, command, grammar, abs, relTestPath(projectRoot, abs))
	}
}
