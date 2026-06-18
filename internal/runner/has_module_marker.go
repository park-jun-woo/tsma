//ff:func feature=runner type=helper control=iteration dimension=1
//ff:what Reports whether a directory holds a Maven or Gradle build marker file
package runner

import (
	"os"
	"path/filepath"
)

// hasModuleMarker reports whether dir contains a Maven or Gradle build marker:
// pom.xml, build.gradle, or build.gradle.kts. It is the per-directory predicate
// used by NearestModuleRoot during its upward walk.
func hasModuleMarker(dir string) bool {
	for _, marker := range []string{"pom.xml", "build.gradle", "build.gradle.kts"} {
		if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
			return true
		}
	}
	return false
}
