//ff:func feature=runner type=helper control=sequence
//ff:what Identifies the Java build tool (maven or gradle) from marker files in a project root
package runner

import (
	"os"
	"path/filepath"
)

// detectJavaBuildTool reports the Java build tool used by the project at
// projectRoot: "maven" for a pom.xml, "gradle" for build.gradle(.kts), and ""
// when neither marker is found. Maven takes priority when both exist.
func detectJavaBuildTool(projectRoot string) string {
	if _, err := os.Stat(filepath.Join(projectRoot, "pom.xml")); err == nil {
		return "maven"
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "build.gradle")); err == nil {
		return "gradle"
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "build.gradle.kts")); err == nil {
		return "gradle"
	}
	return ""
}
