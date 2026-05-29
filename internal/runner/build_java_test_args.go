//ff:func feature=runner type=helper control=selection
//ff:what Builds test command arguments for the given Java build tool and test class
package runner

// buildJavaTestArgs selects the Maven or Gradle argument builder based on the
// build tool and returns the arguments to run a single test class. An empty
// slice is returned for an unknown build tool. Kept environment-independent for
// unit testing.
func buildJavaTestArgs(buildTool, testClass string) []string {
	switch buildTool {
	case "maven":
		return buildMavenTestArgs(testClass)
	case "gradle":
		return buildGradleTestArgs(testClass)
	default:
		return nil
	}
}
