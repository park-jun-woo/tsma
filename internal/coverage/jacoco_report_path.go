//ff:func feature=coverage type=helper control=selection
//ff:what Resolves the JaCoCo XML report path for the given build tool relative to the module root
package coverage

import "path/filepath"

// jacocoReportPath returns the conventional JaCoCo XML report path for the
// given build tool, joined under moduleRoot:
//   - maven:  target/site/jacoco/jacoco.xml
//   - gradle: build/reports/jacoco/test/jacocoTestReport.xml
//
// moduleRoot is the build-module directory resolved by NearestModuleRoot (it
// equals projectRoot for single-module projects), so this is the single point
// that decides where the report is read from for both single- and multi-module
// layouts. An empty string is returned for an unknown build tool.
func jacocoReportPath(moduleRoot, buildTool string) string {
	switch buildTool {
	case "maven":
		return filepath.Join(moduleRoot, "target", "site", "jacoco", "jacoco.xml")
	case "gradle":
		return filepath.Join(moduleRoot, "build", "reports", "jacoco", "test", "jacocoTestReport.xml")
	default:
		return ""
	}
}
