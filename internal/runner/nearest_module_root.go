//ff:func feature=runner type=helper control=iteration dimension=1
//ff:what Walks up from a project-relative file to the nearest Maven/Gradle module root, or projectRoot
package runner

import (
	"path/filepath"
)

// NearestModuleRoot resolves the build-module root for fileRel by walking up the
// directory tree from fileRel toward projectRoot and returning the first
// directory that holds a build marker (pom.xml for Maven, or build.gradle /
// build.gradle.kts for Gradle). The nearest (innermost) marker wins, which
// matches multi-module conventions. projectRoot is included as a stop boundary;
// when no marker is found at or below it, projectRoot is returned (so a
// single-module project falls back to current behavior, regression-free).
//
// projectRoot is treated as absolute and fileRel as project-relative; the
// returned module root is absolute (joined under projectRoot) and so is the
// projectRoot fallback, keeping callers consistent.
func NearestModuleRoot(projectRoot, fileRel string) string {
	dir := filepath.Dir(filepath.Join(projectRoot, fileRel))
	root := filepath.Clean(projectRoot)

	for {
		if hasModuleMarker(dir) {
			return dir
		}
		if dir == root {
			return root
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without hitting projectRoot; fall back.
			return root
		}
		dir = parent
	}
}
