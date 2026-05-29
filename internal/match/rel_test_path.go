//ff:func feature=match type=helper control=sequence lang=go
//ff:what Returns a project-root-relative slash path for a test file (matcher style)
package match

import "path/filepath"

// relTestPath returns the path of absPath relative to projectRoot, mirroring
// the testRel form produced by GoMatcher.Match. On failure it falls back to the
// base name so the index never stores an absolute path.
func relTestPath(projectRoot, absPath string) string {
	rel, err := filepath.Rel(projectRoot, absPath)
	if err != nil {
		return filepath.Base(absPath)
	}
	return rel
}
