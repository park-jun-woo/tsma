//ff:func feature=match type=helper control=iteration dimension=1 lang=java
//ff:what ingestJavaTestDir: scans one JUnit test directory for *Test.java/*Tests.java files and ingests each into the content index (via ingestJavaTestFile). Extracted from BuildJavaPkgTestIndex so the dir scan and the entry loop stay at separate depths (Q1 ≤2). A missing directory is a no-op (graceful fallback to filename matching upstream).
package match

import (
	"os"
	"path/filepath"
)

// ingestJavaTestDir ingests every JUnit test file directly in dir (relative to
// projectRoot) into idx.
func ingestJavaTestDir(idx *JavaPkgTestIndex, projectRoot, dir, command, grammar string) {
	absDir := filepath.Join(projectRoot, dir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !isJavaTestFile(name) {
			continue
		}
		abs := filepath.Join(absDir, name)
		ingestJavaTestFile(idx, command, grammar, abs, relTestPath(projectRoot, abs))
	}
}
