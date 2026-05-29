//ff:func feature=coverage type=helper control=selection
//ff:what Builds the build-tool arguments that run tests and produce a JaCoCo report
package coverage

// buildJavaCoverageArgs returns the build-tool arguments that run the test
// suite and generate a JaCoCo XML report. Maven uses `test jacoco:report`;
// Gradle uses `test jacocoTestReport`. An empty slice is returned for an
// unknown build tool. Kept environment-independent for unit testing.
func buildJavaCoverageArgs(buildTool string) []string {
	switch buildTool {
	case "maven":
		return []string{"test", "jacoco:report"}
	case "gradle":
		return []string{"test", "jacocoTestReport"}
	default:
		return nil
	}
}
