//ff:func feature=coverage type=helper control=selection
//ff:what Resolves the JaCoCo XML report path for the given build tool relative to the project root
package coverage

import "path/filepath"

// jacocoReportPath returns the conventional JaCoCo XML report path for the
// given build tool, joined under projectRoot:
//   - maven:  target/site/jacoco/jacoco.xml
//   - gradle: build/reports/jacoco/test/jacocoTestReport.xml
// An empty string is returned for an unknown build tool.
func jacocoReportPath(projectRoot, buildTool string) string {
	switch buildTool {
	case "maven":
		return filepath.Join(projectRoot, "target", "site", "jacoco", "jacoco.xml")
	case "gradle":
		return filepath.Join(projectRoot, "build", "reports", "jacoco", "test", "jacocoTestReport.xml")
	default:
		return ""
	}
}
